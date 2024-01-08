import fs from 'fs';
import MonacoThemes from 'monaco-themes';

const file = fs.readFileSync('src/static/onedarkTextmate.tmtheme', {encoding: 'utf-8'});
fs.writeFileSync('src/components/monaco-editor-theme/onedark/OnedarkTextmate.tmtheme.ts', 'export default ' + JSON.stringify(MonacoThemes.parseTmTheme(file)));
