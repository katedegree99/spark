import { defineConfig } from 'orval';

export default defineConfig({
  spark: {
    input: {
      target: '../schema/openapi/openapi.yaml',
    },
    output: {
      client: 'swr',
      target: './src/lib/api/generated.ts',
      schemas: './src/lib/api/model',
      clean: true,
    },
  },
});
