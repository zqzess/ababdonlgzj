/**
 * ========================
 * Created with WebStorm.
 * @Author zqzess
 * @Date 2024/10/09 20:46
 * @File abandonlgzj/index.ts
 * @Version :
 * @Desc :
 * @GitHUb Https://github.com/zqzess
 * ========================
 **/
import { createRouter, createWebHashHistory } from "vue-router";
import Upload from '../views/Upload.vue'
import Home from "../views/Home.vue";

const router = createRouter({
    // history: createWebHistory(import.meta.env.BASE_URL),
    history: createWebHashHistory(),
    routes: [
        {
            path: '/',
            redirect: '/upload',
        },
        {
            path: '/upload',
            name: 'Upload',
            component: Upload,
            // redirect: '/home',
            // children: [
            //     {
            //         path: '/home',
            //         name: 'FileList',
            //         component: FileList,
            //         meta: {
            //             title: '主页',
            //         }
            //     },
            // ]
        },
        {
            path: '/home',
            name: 'Home',
            component: Home
        },
        {
            path: '/:catchAll(.*)',
            component: () => import(/* webpackChunkName: "NotFound" */ '../views/NotFound.vue'),
            meta: {
                title: 'notFound',
            },
        },
    ]
})

router.beforeEach((to, from, next) => {
    document.title = `${to.meta.title} | 灵敢足迹备份工具`;
    const role = sessionStorage.getItem('session');
    if (!role && to.path !== '/upload' && to.path !== '/regist') {
        next('/upload');
    }
    else
    {
        next()
    }
});
export default router
