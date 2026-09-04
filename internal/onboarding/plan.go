package onboarding

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Tom-R-Main/Internationalizer/internal/config"
	"github.com/Tom-R-Main/Internationalizer/internal/message"
	"gopkg.in/yaml.v3"
)

const planVersion = 1
const maxPlanConfigSize = 2 << 20

// PlanOptions contains explicit user decisions; discovery alone never selects a
// new bundle or overwrites an existing bundle's runtime profile.
type PlanOptions struct {
	AddBundles     []config.Bundle
	Syntax         map[string]message.Syntax
	ConfirmSources []string
	SourceLocale   string
	TargetLocales  []string
}

type PlanDecision struct {
	Code    string `json:"code"`
	Bundle  string `json:"bundle,omitempty"`
	Source  string `json:"source,omitempty"`
	Message string `json:"message"`
}

// Observation binds the proposal to the source and runtime evidence it used.
// Missing paths are recorded too, so newly-created evidence invalidates a plan.
type Observation struct {
	Path   string `json:"path"`
	Exists bool   `json:"exists"`
	SHA256 string `json:"sha256,omitempty"`
}

// ConfigPlan is a reviewable, versioned proposal, not an authorization token.
// Its content hash detects accidental modification, not a malicious signer.
type ConfigPlan struct {
	SchemaVersion     int            `json:"schema_version"`
	ID                string         `json:"id"`
	Root              string         `json:"root"`
	ConfigPath        string         `json:"config_path"`
	BeforeExists      bool           `json:"before_exists"`
	BeforeSHA256      string         `json:"before_sha256,omitempty"`
	AfterSHA256       string         `json:"after_sha256"`
	DiscoverySHA256   string         `json:"discovery_sha256"`
	ProposedYAML      string         `json:"proposed_yaml"`
	Diff              string         `json:"diff"`
	Observations      []Observation  `json:"observations"`
	RequiredDecisions []PlanDecision `json:"required_decisions"`
}

type ApplyReceipt struct {
	SchemaVersion           int      `json:"schema_version"`
	PlanID                  string   `json:"plan_id"`
	Status                  string   `json:"status"`
	ChangedPaths            []string `json:"changed_paths"`
	ConfigPath              string   `json:"config_path"`
	ConfigSHA256            string   `json:"config_sha256"`
	ObservationsRevalidated bool     `json:"observations_revalidated"`
}

type PlanError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *PlanError) Error() string    { return e.Message }
func (e *PlanError) JSONCode() string { return e.Code }

func planError(code, message string) error { return &PlanError{Code: code, Message: message} }

// BuildPlan reads project state and returns an immutable proposal. It performs
// no writes, provider calls or credential materialization.
func BuildPlan(root, configPath string, options PlanOptions) (*ConfigPlan, error) {
	root, err := canonicalPlanRoot(root)
	if err != nil {
		return nil, err
	}
	configPath, err = planConfigPath(root, configPath)
	if err != nil {
		return nil, err
	}
	before, exists, err := readPlanFile(root, configPath)
	if err != nil {
		return nil, err
	}
	if len(before) > maxPlanConfigSize {
		return nil, planError("config_too_large", "configuration exceeds the plan size limit")
	}
	doc := &yaml.Node{Kind: yaml.DocumentNode, Content: []*yaml.Node{{Kind: yaml.MappingNode, Tag: "!!map"}}}
	if exists {
		if err := yaml.Unmarshal(before, doc); err != nil {
			return nil, planError("invalid_config", "configuration is not valid YAML")
		}
	}
	if len(doc.Content) != 1 || doc.Content[0].Kind != yaml.MappingNode {
		return nil, planError("invalid_config", "configuration must be a YAML mapping")
	}
	if err := inspectPlanYAML(doc); err != nil {
		return nil, err
	}
	var cfg config.Config
	if err := doc.Decode(&cfg); err != nil {
		return nil, planError("invalid_config", "configuration fields have invalid types")
	}
	cfg.ApplyDefaults()
	inspection, err := Scan(root, configPath)
	if err != nil {
		return nil, err
	}
	content := doc.Content[0]
	if options.SourceLocale != "" {
		setPlanScalar(content, "source_locale", options.SourceLocale)
	}
	if len(options.TargetLocales) > 0 {
		values := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		for _, locale := range options.TargetLocales {
			values.Content = append(values.Content, planScalar(locale))
		}
		setPlanNode(content, "target_locales", values)
	}
	if !exists && options.SourceLocale == "" {
		setPlanScalar(content, "source_locale", "en")
	}
	if len(options.AddBundles) > 0 || len(options.Syntax) > 0 {
		if len(cfg.Bundles) == 0 && cfg.SourcePath != "" {
			bundles := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
			for _, b := range cfg.EffectiveBundles() {
				node := planBundleNode(b)
				// Carry source_path comments onto its explicit bundle equivalent.
				if source := getPlanNode(content, "source_path"); source != nil {
					copy := *source
					setPlanNode(node, "source", &copy)
					getPlanNode(node, "source").HeadComment = source.HeadComment
					getPlanNode(node, "source").LineComment = source.LineComment
					getPlanNode(node, "source").FootComment = source.FootComment
					if key := getPlanKeyNode(content, "source_path"); key != nil {
						newKey := getPlanKeyNode(node, "source")
						newKey.HeadComment, newKey.LineComment, newKey.FootComment = key.HeadComment, key.LineComment, key.FootComment
					}
				}
				bundles.Content = append(bundles.Content, node)
			}
			setPlanNode(content, "bundles", bundles)
			removePlanKey(content, "source_path")
		}
	}
	bundleNodes := getPlanNode(content, "bundles")
	if bundleNodes == nil && len(options.AddBundles) > 0 {
		bundleNodes = &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
		setPlanNode(content, "bundles", bundleNodes)
	}
	seen := map[string]*yaml.Node{}
	if bundleNodes != nil {
		if bundleNodes.Kind != yaml.SequenceNode {
			return nil, planError("invalid_config", "bundles must be a sequence")
		}
		for _, b := range bundleNodes.Content {
			id := getPlanNode(b, "id")
			if id == nil || id.Value == "" || seen[id.Value] != nil {
				return nil, planError("invalid_config", "existing bundle identities must be nonempty and unique")
			}
			seen[id.Value] = b
		}
	}
	additions := append([]config.Bundle(nil), options.AddBundles...)
	sort.Slice(additions, func(i, j int) bool { return additions[i].ID < additions[j].ID })
	for _, b := range additions {
		if b.ID == "" || seen[b.ID] != nil {
			return nil, planError("invalid_decision", "added bundle identity must be nonempty and not already configured")
		}
		if b.Source == "" || b.Target == "" {
			return nil, planError("invalid_decision", "added bundles require explicit source and target paths")
		}
		if _, err := safePlanPath(root, b.Source); err != nil {
			return nil, err
		}
		if _, err := safePlanPath(root, b.Target); err != nil {
			return nil, err
		}
		node := planBundleNode(b)
		bundleNodes.Content = append(bundleNodes.Content, node)
		seen[b.ID] = node
	}
	for id, syntax := range options.Syntax {
		if seen[id] == nil {
			return nil, planError("invalid_decision", "syntax selection references an unknown bundle: "+id)
		}
		if err := message.ValidateSyntax(syntax); err != nil || syntax == "" {
			return nil, planError("invalid_decision", "syntax selection must be auto, plain, i18next, or icu")
		}
		setPlanScalar(seen[id], "message_syntax", string(syntax))
	}
	var proposed config.Config
	if err := doc.Decode(&proposed); err != nil {
		return nil, planError("invalid_config", "proposed configuration has invalid field types")
	}
	proposed.ApplyDefaults()
	plan := &ConfigPlan{SchemaVersion: planVersion, Root: root, ConfigPath: configPath, BeforeExists: exists, Observations: []Observation{}, RequiredDecisions: []PlanDecision{}}
	plan.DiscoverySHA256, err = discoveryPlanDigest(inspection)
	if err != nil {
		return nil, err
	}
	if exists {
		plan.BeforeSHA256 = planHash(before)
	}
	if len(proposed.TargetLocales) == 0 {
		plan.RequiredDecisions = append(plan.RequiredDecisions, PlanDecision{Code: "TARGET_LOCALES_REQUIRED", Message: "Select target locales before applying the configuration."})
	}
	if len(proposed.EffectiveBundles()) == 0 {
		plan.RequiredDecisions = append(plan.RequiredDecisions, PlanDecision{Code: "BUNDLE_SELECTION_REQUIRED", Message: "Select an authoritative source catalog and target template."})
	}
	if len(proposed.TargetLocales) > 0 && len(proposed.EffectiveBundles()) > 0 {
		if err := proposed.ValidateProject(); err != nil {
			return nil, planError("invalid_config", "proposed configuration: "+err.Error())
		}
	}
	confirmed := map[string]bool{}
	for _, source := range options.ConfirmSources {
		path, err := safePlanPath(root, source)
		if err != nil {
			return nil, err
		}
		confirmed[path] = true
	}
	paths := map[string]bool{}
	for _, b := range proposed.EffectiveBundles() {
		source, err := safePlanPath(root, b.Source)
		if err != nil {
			return nil, err
		}
		paths[source] = true
		if _, exists, err := readPlanFile(root, source); err != nil {
			return nil, err
		} else if !exists {
			plan.RequiredDecisions = append(plan.RequiredDecisions, PlanDecision{Code: "SOURCE_NOT_FOUND", Bundle: b.ID, Source: b.Source, Message: "The selected source catalog does not exist; choose an existing authoritative path."})
		}
		for _, locale := range proposed.TargetLocales {
			target, err := b.TargetPath(locale)
			if err != nil {
				return nil, planError("invalid_config", err.Error())
			}
			if _, err := safePlanPath(root, target); err != nil {
				return nil, err
			}
		}
		if temporaryPlanPath(root, source) && !confirmed[source] {
			plan.RequiredDecisions = append(plan.RequiredDecisions, PlanDecision{Code: "SOURCE_CONFIRMATION_REQUIRED", Bundle: b.ID, Source: b.Source, Message: "Confirm that this temporary-path catalog is the authoritative source; temporary paths are not automatically invalid."})
		}
		for _, d := range inspection.Diagnostics {
			if d.Code == "AUTO_SYNTAX_AMBIGUOUS" && d.Bundle == b.ID && b.MessageSyntax == message.Auto {
				plan.RequiredDecisions = append(plan.RequiredDecisions, PlanDecision{Code: "SYNTAX_SELECTION_REQUIRED", Bundle: b.ID, Source: b.Source, Message: "Select plain, i18next, or icu for ambiguous automatic message syntax."})
				break
			}
		}
	}
	for _, b := range additions {
		selected := b.MessageSyntax
		if syntax, ok := options.Syntax[b.ID]; ok {
			selected = syntax
		}
		intrinsicFluent := strings.EqualFold(b.Format, "fluent") || strings.EqualFold(filepath.Ext(b.Source), ".ftl")
		if !intrinsicFluent && (selected == "" || selected == message.Auto) {
			plan.RequiredDecisions = append(plan.RequiredDecisions, PlanDecision{Code: "SYNTAX_SELECTION_REQUIRED", Bundle: b.ID, Source: b.Source, Message: "Explicitly choose the new bundle's runtime syntax from the discovery evidence."})
		}
	}
	for _, candidate := range inspection.Candidates {
		paths[candidate.Source] = true
		for _, evidence := range candidate.Evidence {
			if evidence.Path != "" {
				paths[evidence.Path] = true
			}
		}
	}
	for _, b := range inspection.Bundles {
		for _, evidence := range b.Evidence {
			if evidence.Path != "" {
				paths[evidence.Path] = true
			}
		}
	}
	// Runtime integration may be installed after discovery. Record absent package
	// manifests between each catalog and root as well as the evidence we found.
	for _, b := range proposed.EffectiveBundles() {
		source, _ := safePlanPath(root, b.Source)
		for dir := filepath.Dir(source); ; dir = filepath.Dir(dir) {
			paths[filepath.Join(dir, "package.json")] = true
			if dir == root {
				break
			}
		}
	}
	for path := range paths {
		abs, err := safePlanPath(root, path)
		if err != nil {
			return nil, err
		}
		if abs == configPath {
			continue
		}
		data, present, err := readPlanFile(root, abs)
		if err != nil {
			return nil, err
		}
		ob := Observation{Path: abs, Exists: present}
		if present {
			ob.SHA256 = planHash(data)
		}
		plan.Observations = append(plan.Observations, ob)
	}
	sort.Slice(plan.Observations, func(i, j int) bool { return plan.Observations[i].Path < plan.Observations[j].Path })
	sort.Slice(plan.RequiredDecisions, func(i, j int) bool {
		a, b := plan.RequiredDecisions[i], plan.RequiredDecisions[j]
		return a.Bundle+"\x00"+a.Code < b.Bundle+"\x00"+b.Code
	})
	var encoded bytes.Buffer
	encoder := yaml.NewEncoder(&encoded)
	encoder.SetIndent(2)
	if err := encoder.Encode(doc); err != nil {
		return nil, err
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	plan.ProposedYAML = encoded.String()
	// No-op planning must not rewrite formatting or comments.
	if exists && len(options.AddBundles) == 0 && len(options.Syntax) == 0 && options.SourceLocale == "" && len(options.TargetLocales) == 0 {
		plan.ProposedYAML = string(before)
	}
	plan.AfterSHA256 = planHash([]byte(plan.ProposedYAML))
	plan.Diff = planDiff(configPath, string(before), plan.ProposedYAML)
	plan.ID, err = planDigest(plan)
	return plan, err
}

// ApplyPlan rechecks all observations under a cooperative per-config lock and
// replaces only the exact config file. It never executes commands from a plan.
func ApplyPlan(plan *ConfigPlan) (*ApplyReceipt, error) {
	if plan == nil || plan.SchemaVersion != planVersion {
		return nil, planError("invalid_plan", "unsupported or missing configuration plan")
	}
	digest, err := planDigest(plan)
	if err != nil || plan.ID != digest || plan.AfterSHA256 != planHash([]byte(plan.ProposedYAML)) {
		return nil, planError("invalid_plan", "plan integrity check failed; create a new plan")
	}
	if len(plan.RequiredDecisions) > 0 {
		return nil, planError("decisions_required", "plan has unresolved decisions; create a new plan with explicit selections")
	}
	if len(plan.ProposedYAML) > maxPlanConfigSize {
		return nil, planError("invalid_plan", "proposed configuration exceeds the plan size limit")
	}
	root, err := canonicalPlanRoot(plan.Root)
	if err != nil {
		return nil, err
	}
	path, err := safePlanPath(root, plan.ConfigPath)
	if err != nil {
		return nil, err
	}
	var doc yaml.Node
	if err := yaml.Unmarshal([]byte(plan.ProposedYAML), &doc); err != nil {
		return nil, planError("invalid_plan", "proposed configuration is not valid YAML")
	}
	if err := inspectPlanYAML(&doc); err != nil {
		return nil, err
	}
	var cfg config.Config
	if err := doc.Decode(&cfg); err != nil {
		return nil, planError("invalid_plan", "proposed configuration fields have invalid types")
	}
	cfg.ApplyDefaults()
	if err := cfg.ValidateProject(); err != nil {
		return nil, planError("invalid_plan", "proposed configuration is invalid: "+err.Error())
	}
	for _, b := range cfg.EffectiveBundles() {
		if _, err := safePlanPath(root, b.Source); err != nil {
			return nil, err
		}
		for _, locale := range cfg.TargetLocales {
			target, _ := b.TargetPath(locale)
			if _, err := safePlanPath(root, target); err != nil {
				return nil, err
			}
		}
	}
	lockPath := path + ".apply-lock"
	if _, err := safePlanPath(root, lockPath); err != nil {
		return nil, err
	}
	lock, err := os.OpenFile(lockPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if errors.Is(err, os.ErrExist) {
		return nil, planError("apply_locked", "another config application owns the lock; inspect it before retrying")
	}
	if err != nil {
		return nil, planError("apply_io", "cannot create the configuration application lock")
	}
	defer func() { _ = lock.Close(); _ = os.Remove(lockPath) }()
	before, exists, err := readPlanFile(root, path)
	if err != nil {
		return nil, err
	}
	receipt := &ApplyReceipt{SchemaVersion: planVersion, PlanID: plan.ID, ConfigPath: path, ConfigSHA256: plan.AfterSHA256, ChangedPaths: []string{}}
	// Replay attests only that config bytes already match the plan. Source and
	// framework evidence can evolve after a successful application; callers must
	// run config check/dry-run to establish current translation readiness.
	if exists && planHash(before) == plan.AfterSHA256 {
		receipt.Status = "already_applied"
		return receipt, nil
	}
	if exists != plan.BeforeExists || (exists && planHash(before) != plan.BeforeSHA256) {
		return nil, planError("stale_plan", "configuration changed; create a new configuration plan")
	}
	for _, observation := range plan.Observations {
		data, exists, err := readPlanFile(root, observation.Path)
		if err != nil {
			return nil, err
		}
		if exists != observation.Exists || (exists && planHash(data) != observation.SHA256) {
			return nil, planError("stale_plan", "source or runtime evidence changed; create a new configuration plan")
		}
	}
	inspection, err := Scan(root, path)
	if err != nil {
		return nil, planError("stale_plan", "project discovery can no longer reproduce the saved plan; create a new plan")
	}
	discoveryDigest, err := discoveryPlanDigest(inspection)
	if err != nil || discoveryDigest != plan.DiscoverySHA256 {
		return nil, planError("stale_plan", "discovered catalogs or runtime evidence changed; create a new configuration plan")
	}
	receipt.ObservationsRevalidated = true
	mode := os.FileMode(0600)
	if exists {
		info, err := os.Stat(path)
		if err != nil {
			return nil, planError("apply_io", "cannot inspect configuration permissions")
		}
		mode = info.Mode().Perm()
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".internationalizer-apply-*")
	if err != nil {
		return nil, planError("apply_io", "cannot prepare the configuration replacement")
	}
	tmpPath := tmp.Name()
	defer func() { _ = tmp.Close(); _ = os.Remove(tmpPath) }()
	if err := tmp.Chmod(mode); err != nil {
		return nil, planError("apply_io", "cannot preserve configuration permissions")
	}
	if _, err := tmp.WriteString(plan.ProposedYAML); err != nil {
		return nil, planError("apply_io", "cannot write the configuration replacement")
	}
	if err := tmp.Sync(); err != nil {
		return nil, planError("apply_io", "cannot flush the configuration replacement")
	}
	if err := tmp.Close(); err != nil {
		return nil, planError("apply_io", "cannot close the configuration replacement")
	}
	// Recheck immediately before commit as protection against ordinary editors
	// that do not participate in the cooperative lock.
	current, currentExists, err := readPlanFile(root, path)
	if err != nil {
		return nil, err
	}
	if currentExists != exists || !bytes.Equal(current, before) {
		return nil, planError("stale_plan", "configuration changed during application; create a new plan")
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return nil, planError("apply_io", "cannot commit the configuration replacement")
	}
	readback, present, err := readPlanFile(root, path)
	if err != nil || !present || planHash(readback) != plan.AfterSHA256 {
		return nil, planError("apply_verification_failed", "configuration was written but readback did not match; inspect current state before retrying")
	}
	receipt.Status = "applied"
	receipt.ChangedPaths = []string{path}
	return receipt, nil
}

func planDigest(plan *ConfigPlan) (string, error) {
	copy := *plan
	copy.ID = ""
	data, err := json.Marshal(copy)
	if err != nil {
		return "", err
	}
	return planHash(data), nil
}

func discoveryPlanDigest(inspection *Inspection) (string, error) {
	// Credential presence is neither source evidence nor provider verification.
	// Excluding it keeps otherwise identical offline plans deterministic.
	copy := *inspection
	copy.Credentials = nil
	data, err := json.Marshal(copy)
	if err != nil {
		return "", err
	}
	return planHash(data), nil
}

func planHash(data []byte) string {
	digest := sha256.Sum256(data)
	return hex.EncodeToString(digest[:])
}

func canonicalPlanRoot(root string) (string, error) {
	abs, err := filepath.Abs(root)
	if err != nil {
		return "", planError("unsafe_path", "cannot resolve project root")
	}
	abs, err = filepath.EvalSymlinks(abs)
	if err != nil {
		return "", planError("unsafe_path", "project root must exist")
	}
	info, err := os.Stat(abs)
	if err != nil || !info.IsDir() {
		return "", planError("unsafe_path", "project root must be a directory")
	}
	return abs, nil
}

func planConfigPath(root, path string) (string, error) {
	if path != "" {
		return safePlanPath(root, path)
	}
	resolved := resolveConfigPath(root, "")
	rel, err := filepath.Rel(root, resolved)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", planError("external_config_scope", "the active default configuration is outside the project; use --config with an explicit project-local path to create a local configuration intentionally")
	}
	return safePlanPath(root, resolved)
}

// Reject links in every project-relative component, including dangling links.
// Root itself is canonicalized to support platform aliases such as /var.
func safePlanPath(root, path string) (string, error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, path)
	}
	path = filepath.Clean(path)
	if sensitivePlanPath(path) {
		return "", planError("unsafe_path", "configuration plans do not read or write credential-shaped files")
	}
	rel, err := filepath.Rel(root, path)
	if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", planError("unsafe_path", "configuration plans may only access files inside the project root")
	}
	current := root
	parts := strings.Split(rel, string(filepath.Separator))
	for i, part := range parts {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return "", planError("unsafe_path", "cannot inspect a project path safely")
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return "", planError("unsafe_path", "configuration plans do not follow symlinks")
		}
		if i < len(parts)-1 && !info.IsDir() {
			return "", planError("unsafe_path", "project path has a non-directory ancestor")
		}
		if i == len(parts)-1 && !info.Mode().IsRegular() {
			return "", planError("unsafe_path", "configuration plans require regular files")
		}
	}
	return path, nil
}

func sensitivePlanPath(path string) bool {
	name := strings.ToLower(filepath.Base(path))
	return name == ".env" || strings.HasPrefix(name, ".env.") || strings.HasSuffix(name, ".key") || strings.HasSuffix(name, ".pem") || strings.HasSuffix(name, ".p12")
}

func readPlanFile(root, path string) ([]byte, bool, error) {
	path, err := safePlanPath(root, path)
	if err != nil {
		return nil, false, err
	}
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, planError("apply_io", "cannot read a project file required by the plan")
	}
	defer func() { _ = file.Close() }()
	data, err := io.ReadAll(io.LimitReader(file, maxDiscoveryBytes+1))
	if err != nil {
		return nil, false, planError("apply_io", "cannot read a project file required by the plan")
	}
	if len(data) > maxDiscoveryBytes {
		return nil, false, planError("file_too_large", "a project file exceeds the plan observation size limit")
	}
	return data, true, nil
}

func temporaryPlanPath(root, path string) bool {
	rel, _ := filepath.Rel(root, path)
	for _, part := range strings.Split(filepath.ToSlash(rel), "/") {
		if part == "tmp" || part == "temp" {
			return true
		}
	}
	return false
}

// Plans preserve unknown settings, but cannot safely publish opaque inline
// secrets. Reject such configs before serializing either YAML or a diff.
func inspectPlanYAML(node *yaml.Node) error {
	if node.Kind == yaml.AliasNode {
		return planError("unsupported_config", "configuration plan editing does not support YAML aliases; expand aliases before planning")
	}
	if node.Kind == yaml.MappingNode {
		seen := map[string]bool{}
		for i := 0; i+1 < len(node.Content); i += 2 {
			key := node.Content[i].Value
			if seen[key] {
				return planError("invalid_config", "configuration contains duplicate YAML keys")
			}
			seen[key] = true
			lower := strings.ToLower(key)
			if lower == "<<" {
				return planError("unsupported_config", "configuration plan editing does not support YAML merge keys")
			}
			if !strings.HasSuffix(lower, "_env") && (strings.Contains(lower, "password") || strings.Contains(lower, "secret") || strings.Contains(lower, "token") || lower == "api_key" || strings.HasSuffix(lower, "_api_key") || lower == "private_key") {
				return planError("inline_secret", "configuration contains a possible inline credential; use environment variable references before creating a saved plan")
			}
		}
	}
	if node.Kind == yaml.ScalarNode {
		if parsed, err := url.Parse(node.Value); err == nil && parsed.Scheme != "" && parsed.Host != "" {
			if parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
				return planError("inline_secret", "configuration contains a URL with possible embedded credentials; use a credential-free endpoint before creating a saved plan")
			}
		}
	}
	for _, child := range node.Content {
		if err := inspectPlanYAML(child); err != nil {
			return err
		}
	}
	return nil
}

func planScalar(value string) *yaml.Node {
	return &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!str", Value: value}
}
func getPlanNode(mapping *yaml.Node, key string) *yaml.Node {
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i+1]
		}
	}
	return nil
}
func getPlanKeyNode(mapping *yaml.Node, key string) *yaml.Node {
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			return mapping.Content[i]
		}
	}
	return nil
}
func setPlanNode(mapping *yaml.Node, key string, value *yaml.Node) {
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			old := mapping.Content[i+1]
			value.HeadComment = old.HeadComment
			value.LineComment = old.LineComment
			value.FootComment = old.FootComment
			mapping.Content[i+1] = value
			return
		}
	}
	mapping.Content = append(mapping.Content, planScalar(key), value)
}
func setPlanScalar(mapping *yaml.Node, key, value string) {
	if node := getPlanNode(mapping, key); node != nil {
		node.Kind = yaml.ScalarNode
		node.Tag = "!!str"
		node.Value = value
		node.Content = nil
		return
	}
	setPlanNode(mapping, key, planScalar(value))
}
func removePlanKey(mapping *yaml.Node, key string) {
	for i := 0; i+1 < len(mapping.Content); i += 2 {
		if mapping.Content[i].Value == key {
			mapping.Content = append(mapping.Content[:i], mapping.Content[i+2:]...)
			return
		}
	}
}
func planBundleNode(b config.Bundle) *yaml.Node {
	node := &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map"}
	setPlanScalar(node, "id", b.ID)
	setPlanScalar(node, "source", b.Source)
	setPlanScalar(node, "target", b.Target)
	if b.Format != "" {
		setPlanScalar(node, "format", b.Format)
	}
	if b.MessageSyntax != "" {
		setPlanScalar(node, "message_syntax", string(b.MessageSyntax))
	}
	return node
}

func planDiff(path, before, after string) string {
	if before == after {
		return ""
	}
	var out strings.Builder
	fmt.Fprintf(&out, "--- %s\n+++ %s\n", path, path)
	oldLines := strings.Split(strings.TrimSuffix(before, "\n"), "\n")
	newLines := strings.Split(strings.TrimSuffix(after, "\n"), "\n")
	if before == "" {
		oldLines = nil
	}
	if after == "" {
		newLines = nil
	}
	fmt.Fprintf(&out, "@@ -1,%d +1,%d @@\n", len(oldLines), len(newLines))
	for _, line := range oldLines {
		out.WriteString("-" + line + "\n")
	}
	for _, line := range newLines {
		out.WriteString("+" + line + "\n")
	}
	return out.String()
}
