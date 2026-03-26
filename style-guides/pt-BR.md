# Portuguese (Brazil) (pt-BR) — Translation Style Guide

## Language Profile
- Locale: `pt-BR`
- Script: Latin, LTR
- CLDR plural forms: one, other
- Text expansion vs English: 15-25%

## Tone & Formality
O tom deve ser amigável, claro e profissional. Em interfaces de software (UI), utilize o tratamento "você" (informal, mas respeitoso) e evite "tu" ou "vós", a menos que seja um aplicativo com foco altamente regional. A voz deve ser direta e ativa, transmitindo confiança. Evite jargões excessivamente técnicos quando uma palavra simples for suficiente. O objetivo é guiar o usuário de forma natural, empática e sem atrito.

## Grammar
- **Botões e Rótulos (Infinitivo):** Para botões de ação isolados, use sempre verbos no infinitivo, não no imperativo. Exemplo: use "Salvar", "Editar" e "Excluir" (em vez de "Salve", "Edite").
- **Instruções e Mensagens (Imperativo):** Quando estiver instruindo o usuário em uma frase completa, use o imperativo na 3ª pessoa do singular (concordando com "você"). Exemplo: "Clique aqui para continuar" ou "Insira sua senha".
- **Capitalização (Sentence case):** O português brasileiro utiliza maiúsculas apenas na primeira letra da frase e em nomes próprios. Evite o "Title Case" do inglês. Exemplo: traduza "Account Settings" como "Configurações da conta" (e não "Configurações Da Conta").
- **Voz Ativa:** Prefira sempre a voz ativa para tornar as mensagens de erro e sucesso mais diretas. Em vez de "O arquivo foi salvo com sucesso", use "Arquivo salvo" ou "Você salvou o arquivo".
- **Gênero Neutro:** Sempre que possível, reescreva frases para evitar a marcação de gênero se o gênero do usuário for desconhecido. Exemplo: em vez de "Bem-vindo" ou "Bem-vinda", use "Boas-vindas". Nunca utilize "x" ou "@" (ex: "Bem-vindx") em textos de UI profissionais.

## Pluralization
O português brasileiro utiliza principalmente duas categorias no sistema CLDR para pluralização:
- **One:** Usado estritamente para o número 1. Exemplo: "1 arquivo", "1 erro".
- **Other:** Usado para o número 0 e para números de 2 em diante. Exemplo: "0 arquivos", "2 arquivos", "10 arquivos". (Nota: embora na linguagem falada o zero às vezes leve o singular, o padrão técnico e gramatical para UI no Brasil é tratar o zero no plural).

## Punctuation & Typography
- **Aspas:** Use aspas duplas curvas (“ ”) para citações principais e aspas simples (‘ ’) para citações dentro de citações.
- **Números:** Use vírgula (,) como separador decimal e ponto (.) como separador de milhares. Exemplo: 1.234,56.
- **Data:** O formato padrão é numérico na ordem Dia/Mês/Ano (DD/MM/AAAA). Exemplo: 31/12/2023.
- **Hora:** Utilize o formato de 24 horas com dois pontos (ex: 14:30) ou com a letra "h" indicando as horas (ex: 14h30). Evite o uso de AM/PM.
- **Pontuação Final:** Não utilize ponto final em botões, rótulos (labels), títulos ou itens curtos de listas. Reserve o ponto final apenas para frases completas (como descrições e mensagens de erro).

## Terminology

| English | Portuguese (Brazil) | Notes |
|---------|---------------------|-------|
| Save | Salvar | Verbo no infinitivo, padrão para botões de ação. |
| Cancel | Cancelar | Verbo no infinitivo. |
| Delete | Excluir | "Excluir" é o termo preferido em UI. Evite "Deletar" (anglicismo) ou "Apagar". |
| Settings | Configurações | Sempre no plural. Evite "Ajustes" a menos que seja padrão do SO (ex: iOS). |
| Search | Pesquisar | "Pesquisar" ou "Buscar" são aceitos, mas "Pesquisar" é mais comum em campos de texto. |
| Error | Erro | |
| Loading | Carregando | O gerúndio é o padrão no Brasil para indicar uma ação em andamento. |
| Dashboard | Painel | "Painel" ou "Painel de controle". "Dashboard" é aceito apenas em contextos muito técnicos. |
| Notifications | Notificações | |
| Sign in | Entrar | "Fazer login" também é comum, mas "Entrar" é mais curto e direto para botões. |
| Sign out | Sair | "Fazer logout" é aceito, mas "Sair" é a melhor prática para UI limpa. |
| Submit | Enviar | Falso cognato: nunca traduza como "Submeter" em formulários. Use "Enviar". |
| Profile | Perfil | |
| Help | Ajuda | |
| Close | Fechar | Verbo no infinitivo. |
