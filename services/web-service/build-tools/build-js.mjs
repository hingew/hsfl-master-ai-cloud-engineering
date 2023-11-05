import esbuild from 'esbuild';

esbuild.build({
    entryPoints: ['client/src/index.ts'],
    outfile: 'public/js/bundle.js',
    bundle: true,
});