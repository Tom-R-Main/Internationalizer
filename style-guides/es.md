# Spanish (es) — Translation Style Guide

## Language Profile
- Locale: `es`
- Script: Latin, LTR
- CLDR plural forms: one, other
- Text expansion vs English: 20-30%

## Tone & Formality
El tono debe ser profesional, claro, conciso y accesible. Para la interfaz de usuario (UI) y la documentación de software, utiliza el trato formal ("usted") para dirigirte al usuario. Este es el estándar más seguro y universal para el español internacional o neutro, garantizando respeto y claridad. Sé directo pero cortés, evitando la jerga local, los regionalismos y el lenguaje excesivamente coloquial. La voz debe ser activa y orientada a guiar al usuario de manera eficiente.

## Grammar
- **Botones y llamadas a la acción (CTAs):** Usa siempre el infinitivo para las acciones que el usuario va a realizar (ej. "Guardar", "Eliminar", "Copiar"). No uses el imperativo ("Guarde", "Guarda").
- **Uso de mayúsculas (Capitalization):** A diferencia del inglés (Title Case), en español solo se escribe con mayúscula inicial la primera palabra de títulos, encabezados, menús y botones (Sentence case), además de los nombres propios (ej. "Configuración de la cuenta", no "Configuración De La Cuenta").
- **Voz pasiva vs. activa:** El inglés abusa de la voz pasiva. En español, prefiere la voz activa o la pasiva refleja con "se" (ej. "Se produjo un error" en lugar de "Un error fue producido").
- **Gerundios:** Evita la traducción literal del "-ing" inglés como gerundio ("-ando", "-iendo") en títulos o encabezados. Usa sustantivos o infinitivos (ej. "Carga de archivos" en lugar de "Cargando archivos"). El gerundio solo es aceptable para indicar una acción en progreso (ej. "Cargando...").

## Pluralization
El español utiliza dos categorías plurales según el estándar CLDR:
- **One (Uno):** Se utiliza exclusivamente para el número 1 (ej. "1 archivo", "1 error").
- **Other (Otros):** Se utiliza para el cero y todos los demás números (ej. "0 archivos", "2 archivos", "10 errores").
Asegúrate siempre de que los artículos y adjetivos concuerden en género y número con el sustantivo (ej. "Los archivos seleccionados").

## Punctuation & Typography
- **Signos de interrogación y exclamación:** Es estrictamente obligatorio usar los signos de apertura (¡ !, ¿ ?) en oraciones completas (ej. "¿Estás seguro de que deseas salir?").
- **Números:** En español internacional, se recomienda usar el espacio duro o el punto (.) como separador de miles, y la coma (,) para los decimales (ej. 1.234,56 o 1 234,56).
- **Fechas:** El formato estándar es DD/MM/AAAA (ej. 31/12/2023). Nunca uses el formato de EE. UU. (MM/DD/AAAA).
- **Hora:** Prefiere el formato de 24 horas (ej. 14:30). Si usas el formato de 12 horas, utiliza "a. m." y "p. m." (con minúsculas, espacios y puntos).
- **Comillas:** En UI estándar, las comillas dobles (" ") son aceptables, pero en documentación formal prefiere las comillas angulares (« »).

## Terminology

| English | Spanish | Notes |
|---------|---------|-------|
| Save | Guardar | Usar infinitivo para botones de acción. |
| Cancel | Cancelar | Usar infinitivo. |
| Delete | Eliminar | Preferido sobre "Borrar" en contextos de software. |
| Settings | Configuración | Preferido sobre "Ajustes" para español neutro. |
| Search | Buscar | Usar infinitivo para botones o barras de búsqueda. |
| Error | Error | Mantener igual. |
| Loading | Cargando... | Uso aceptado del gerundio para indicar un proceso activo. |
| Dashboard | Panel | O "Panel de control". Evitar el anglicismo "Dashboard". |
| Notifications | Notificaciones | |
| Sign in | Iniciar sesión | Preferido sobre "Ingresar" o "Conectarse". |
| Sign out | Cerrar sesión | Preferido sobre "Salir" o "Desconectarse". |
| Submit | Enviar | |
| Profile | Perfil | |
| Help | Ayuda | |
| Close | Cerrar | Usar infinitivo. |
