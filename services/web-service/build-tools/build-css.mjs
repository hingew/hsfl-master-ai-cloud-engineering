import esbuild from 'esbuild';
import { sassPlugin } from 'esbuild-sass-plugin';
import { glob } from 'glob';

esbuild.build({
    stdin: {
        contents: glob.sync("**/*.scss").map(f => `@import './${f}';`).join('\n'),
        loader: 'css',
        resolveDir: process.cwd(),
    },
    plugins: [sassPlugin()],
    outfile: 'public/css/bundle.css',
    bundle: true,
    minify: true,
});