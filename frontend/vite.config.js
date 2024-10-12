import { defineConfig } from 'vite';
import { fileURLToPath, URL } from 'node:url';
import AutoImport from 'unplugin-auto-import/vite';
import Components from 'unplugin-vue-components/vite';
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers';
import path from 'path';
import vue from '@vitejs/plugin-vue';
var pathSrc = path.resolve(__dirname, 'src');
// https://vitejs.dev/config/
export default defineConfig({
    plugins: [vue(),
        AutoImport({
            resolvers: [ElementPlusResolver()],
            dts: path.resolve(pathSrc, 'auto-imports.d.ts'),
        }),
        Components({
            resolvers: [ElementPlusResolver()],
            dts: path.resolve(pathSrc, 'components.d.ts')
        }),],
    resolve: {
        alias: {
            '@': fileURLToPath(new URL('./src', import.meta.url))
        }
    },
    server: {
        host: '0.0.0.0', //通过ip的形式访问
        port: 5175, //端口号
        open: true, //自动打开浏览器
        // 配置代理，但是我不推荐前端去代理，
        // 因为打包后便不会在有代理，上线后是个坑
        // proxy: {
        //   '/api': {
        //     target: 'http://API网关所在域名',
        //     changeOrigin: true,
        //     rewrite: (path) => path.replace(/^\/api/, '')
        //   },
        // }
    },
});
