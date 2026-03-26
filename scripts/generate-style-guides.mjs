#!/usr/bin/env node

/**
 * Generate per-language style guides using Gemini via AI Studio.
 * Each guide is written IN the target language (except section headers).
 *
 * Usage:
 *   GOOGLE_AI_STUDIO_API_KEY=... node scripts/generate-style-guides.mjs
 *   GOOGLE_AI_STUDIO_API_KEY=... node scripts/generate-style-guides.mjs --locale ko
 */

import { writeFile, mkdir } from "node:fs/promises";
import { join } from "node:path";

const MODEL = process.env.GEMINI_MODEL || "gemini-3.1-pro-preview";
const API_KEY = process.env.GOOGLE_AI_STUDIO_API_KEY;
if (!API_KEY) {
  console.error("GOOGLE_AI_STUDIO_API_KEY is required");
  process.exit(1);
}

const LOCALES = {
  // Original 17 (ExecuFunction languages)
  ar:      { name: "Arabic",               native: "العربية",      script: "Arabic",    dir: "RTL" },
  bn:      { name: "Bengali",              native: "বাংলা",        script: "Bengali",   dir: "LTR" },
  de:      { name: "German",               native: "Deutsch",      script: "Latin",     dir: "LTR" },
  es:      { name: "Spanish",              native: "Español",      script: "Latin",     dir: "LTR" },
  fr:      { name: "French",               native: "Français",     script: "Latin",     dir: "LTR" },
  hi:      { name: "Hindi",                native: "हिन्दी",         script: "Devanagari",dir: "LTR" },
  id:      { name: "Indonesian",           native: "Indonesia",    script: "Latin",     dir: "LTR" },
  ja:      { name: "Japanese",             native: "日本語",        script: "CJK",       dir: "LTR" },
  pa:      { name: "Punjabi",              native: "ਪੰਜਾਬੀ",       script: "Gurmukhi",  dir: "LTR" },
  "pt-BR": { name: "Portuguese (Brazil)",  native: "Português",    script: "Latin",     dir: "LTR" },
  ru:      { name: "Russian",              native: "Русский",      script: "Cyrillic",  dir: "LTR" },
  te:      { name: "Telugu",               native: "తెలుగు",       script: "Telugu",    dir: "LTR" },
  th:      { name: "Thai",                 native: "ไทย",          script: "Thai",      dir: "LTR" },
  uk:      { name: "Ukrainian",            native: "Українська",   script: "Cyrillic",  dir: "LTR" },
  yue:     { name: "Cantonese",            native: "粵語",          script: "CJK",       dir: "LTR" },
  "zh-CN": { name: "Simplified Chinese",   native: "简体中文",      script: "CJK",       dir: "LTR" },
  "zh-TW": { name: "Traditional Chinese",  native: "繁體中文",      script: "CJK",       dir: "LTR" },
  // New 14 (expanded coverage)
  ko:      { name: "Korean",               native: "한국어",        script: "Hangul",    dir: "LTR" },
  tr:      { name: "Turkish",              native: "Türkçe",       script: "Latin",     dir: "LTR" },
  vi:      { name: "Vietnamese",           native: "Tiếng Việt",   script: "Latin",     dir: "LTR" },
  pl:      { name: "Polish",               native: "Polski",       script: "Latin",     dir: "LTR" },
  nl:      { name: "Dutch",                native: "Nederlands",   script: "Latin",     dir: "LTR" },
  sv:      { name: "Swedish",              native: "Svenska",      script: "Latin",     dir: "LTR" },
  it:      { name: "Italian",              native: "Italiano",     script: "Latin",     dir: "LTR" },
  cs:      { name: "Czech",                native: "Čeština",      script: "Latin",     dir: "LTR" },
  el:      { name: "Greek",                native: "Ελληνικά",     script: "Greek",     dir: "LTR" },
  he:      { name: "Hebrew",               native: "עברית",        script: "Hebrew",    dir: "RTL" },
  ms:      { name: "Malay",                native: "Bahasa Melayu",script: "Latin",     dir: "LTR" },
  fi:      { name: "Finnish",              native: "Suomi",        script: "Latin",     dir: "LTR" },
  da:      { name: "Danish",               native: "Dansk",        script: "Latin",     dir: "LTR" },
  ro:      { name: "Romanian",             native: "Română",       script: "Latin",     dir: "LTR" },
};

// CLDR plural forms by language
const CLDR_PLURALS = {
  ar: "zero, one, two, few, many, other",
  bn: "one, other",
  cs: "one, few, many, other",
  da: "one, other",
  de: "one, other",
  el: "one, other",
  es: "one, many, other",
  fi: "one, other",
  fr: "one, many, other",
  he: "one, two, many, other",
  hi: "one, other",
  id: "other",
  it: "one, many, other",
  ja: "other",
  ko: "other",
  ms: "other",
  nl: "one, other",
  pa: "one, other",
  pl: "one, few, many, other",
  "pt-BR": "one, many, other",
  ro: "one, few, other",
  ru: "one, few, many, other",
  sv: "one, other",
  te: "one, other",
  th: "other",
  tr: "one, other",
  uk: "one, few, many, other",
  vi: "other",
  yue: "other",
  "zh-CN": "other",
  "zh-TW": "other",
};

function buildPrompt(locale, info) {
  const plurals = CLDR_PLURALS[locale] || "one, other";

  return `You are an expert linguist and software localization specialist for ${info.name} (${info.native}).

Write a translation style guide for the locale code "${locale}". This guide will be injected into LLM prompts when translating software UI strings and documentation into ${info.name}.

CRITICAL RULES:
- Write ALL explanatory text, examples, grammar rules, and tone descriptions IN ${info.name} (${info.native}).
- Section headers (##) must stay in English so the file structure is scannable by non-speakers.
- The terminology table column headers stay in English, but the "Notes" column content should be in ${info.name}.
- Keep the file under 3KB. Be concise but comprehensive.
- Do NOT include any markdown code fences or wrapping — output raw markdown only.
- Do NOT start with \`\`\`markdown or end with \`\`\`.

Use this exact structure:

# ${info.name} (${locale}) — Translation Style Guide

## Language Profile
- Locale: \`${locale}\`
- Script: ${info.script}, ${info.dir}
- CLDR plural forms: ${plurals}
- Text expansion vs English: [estimate %]

## Tone & Formality
[Explain the formality register to use for software UI. Should it be formal, informal, polite? Why?
What is the voice — direct, friendly, professional? Write this in ${info.name}.]

## Grammar
[3-5 key grammar rules that LLMs commonly get wrong when translating into ${info.name}.
Focus on rules specific to software/UI text: imperative forms for buttons, noun phrases for labels,
sentence structure for error messages. Write in ${info.name}.]

## Pluralization
[Explain the CLDR plural categories for ${info.name} with concrete numeric examples.
Which numbers map to which category? Write in ${info.name}.]

## Punctuation & Typography
[Quotation mark style, decimal/thousands separators, date format (DD/MM/YYYY vs YYYY-MM-DD etc.),
time format (12h vs 24h), any special punctuation rules. Write in ${info.name}.]

## Terminology
[A table of 12-15 common software terms with approved ${info.name} translations.
Include: Save, Cancel, Delete, Settings, Search, Error, Loading, Dashboard, Notifications,
Sign in, Sign out, Submit, Profile, Help, Close.
The table format must be:]

| English | ${info.name} | Notes |
|---------|${"-".repeat(Math.max(info.name.length, 7))}--|-------|
[Fill in all rows. Notes column in ${info.name} explaining usage context if needed.]`;
}

async function generateGuide(locale, info) {
  const prompt = buildPrompt(locale, info);

  const body = {
    contents: [{ role: "user", parts: [{ text: prompt }] }],
    generationConfig: {
      temperature: 0.3,
      maxOutputTokens: 4096,
    },
  };

  const url = `https://generativelanguage.googleapis.com/v1beta/models/${MODEL}:generateContent?key=${API_KEY}`;

  for (let attempt = 0; attempt < 3; attempt++) {
    try {
      const response = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });

      if (response.status === 429 || response.status >= 500) {
        const backoff = Math.pow(2, attempt) * 2000;
        console.error(`  [${locale}] Rate limited (${response.status}), retrying in ${backoff}ms...`);
        await new Promise((r) => setTimeout(r, backoff));
        continue;
      }

      if (!response.ok) {
        const text = await response.text();
        throw new Error(`API ${response.status}: ${text.slice(0, 200)}`);
      }

      const data = await response.json();
      let text = data.candidates?.[0]?.content?.parts?.[0]?.text;
      if (!text) throw new Error("No text in response");

      // Strip markdown code fences if the model wraps them anyway
      text = text.replace(/^```(?:markdown)?\n?/, "").replace(/\n?```$/, "").trim();

      const tokens = data.usageMetadata;
      console.error(
        `  [${locale}] ${info.native} — ${text.length} chars, ${tokens?.promptTokenCount || "?"}/${tokens?.candidatesTokenCount || "?"} tokens`
      );

      return text;
    } catch (err) {
      if (attempt === 2) throw err;
      console.error(`  [${locale}] Error: ${err.message}, retrying...`);
      await new Promise((r) => setTimeout(r, 2000));
    }
  }
}

async function main() {
  const args = process.argv.slice(2);
  let filterLocale = null;
  for (let i = 0; i < args.length; i++) {
    if (args[i] === "--locale" && args[i + 1]) {
      filterLocale = args[i + 1];
    }
  }

  const outDir = join(process.cwd(), "style-guides");
  await mkdir(outDir, { recursive: true });

  const localesToProcess = filterLocale
    ? { [filterLocale]: LOCALES[filterLocale] }
    : LOCALES;

  if (filterLocale && !LOCALES[filterLocale]) {
    console.error(`Unknown locale: ${filterLocale}`);
    process.exit(1);
  }

  const concurrency = 4;
  const entries = Object.entries(localesToProcess);
  let completed = 0;

  console.error(`Generating style guides for ${entries.length} languages (concurrency: ${concurrency})...\n`);

  // Process in batches
  for (let i = 0; i < entries.length; i += concurrency) {
    const batch = entries.slice(i, i + concurrency);
    const results = await Promise.allSettled(
      batch.map(async ([locale, info]) => {
        const guide = await generateGuide(locale, info);
        const outPath = join(outDir, `${locale}.md`);
        await writeFile(outPath, guide + "\n", "utf8");
        completed++;
        return { locale, info };
      })
    );

    for (const result of results) {
      if (result.status === "rejected") {
        console.error(`  FAILED: ${result.reason.message}`);
      }
    }
  }

  console.error(`\nDone. Generated ${completed}/${entries.length} style guides in ${outDir}/`);
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
