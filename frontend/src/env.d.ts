// declare const process = {
//     HOST: 'http://localhost:5173'
// }

export default process;


declare module '*.env' {
    const content: { [key: string]: string };
    export default content;
}

interface ImportMetaEnv {
    VITE_APP_API_URL: string; // 这里声明你使用的环境变量
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}
