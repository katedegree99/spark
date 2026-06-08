import { defineConfig } from 'orval';

export default defineConfig({
  spark: {
    input: {
      target: '../schema/openapi/openapi.yaml',
    },
    output: {
      client: 'swr',
      // 生成物は generated/ 配下に隔離する。orval の clean は target の出力
      // フォルダを丸ごと掃除するため、手書きの fetcher.ts 等を同居させると
      // generate のたびに消える。生成物だけを generated/ に閉じ込めて回避する。
      target: './src/lib/api/generated/client.ts',
      schemas: './src/lib/api/generated/model',
      clean: true,
    },
  },
});
