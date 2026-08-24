import { defineConfig } from 'orval';

export default defineConfig({
    tradingCardAlbum: {
        input: '../docs/swagger.yaml',
        output: {
            mode: 'tags-split',
            target: './app/services/tradingCardAlbumService.ts',
            schemas: 'app/models',
            client: 'fetch',
            override: {
                mutator: {
                    path: './app/services/CustomFetch.ts',
                    name: 'customFetch',
                },
            },
        },
    },
});