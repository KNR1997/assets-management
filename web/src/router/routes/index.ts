const Layout = () => import('@/layout/index.vue')

export const basicRoutes = [
  {
    path: '/',
    redirect: '/workbench', // 默认跳转到首页
    meta: { order: 0 },
  },
  {
    name: 'Workbench-tsdsd',
    path: '/workbench',
    component: Layout,
    children: [
      {
        path: '',
        component: () => import('@/views/workbench/index.vue'),
        name: 'Workbench',
        meta: {
          title: 'Workbench',
          icon: 'icon-park-outline:workbench',
          affix: true,
        },
      },
    ],
    meta: { order: 1 },
  },
  {
    name: 'Login-tsdsd',
    path: '/login',
    component: () => import('@/views/login/index.vue'),
    isHidden: true,
  },
  {
    name: 'Assets-ddd',
    path: '/assets',
    component: Layout,
    children: [
      {
        path: '',
        component: () => import('@/views/asset/index.vue'),
        name: 'Assets',
        meta: {
          title: 'Assets',
          icon: 'fluent:quiz-20-regular',
          affix: true,
        },
      },
    ],
    meta: { order: 2 },
  },
  {
    name: 'Category-ddd',
    path: '/categories',
    component: Layout,
    children: [
      {
        path: '',
        component: () => import('@/views/category/index.vue'),
        name: 'Categories',
        meta: {
          title: 'Categories',
          icon: 'fluent:quiz-20-regular',
          affix: true,
        },
      },
    ],
    meta: { order: 2 },
  },
  {
    name: 'Quizzes-dev',
    path: '/quizzes',
    component: Layout,
    children: [
      {
        path: '',
        component: () => import('@/views/quiz/index.vue'),
        name: 'Quizzes',
        meta: {
          title: 'Quizzes',
          icon: 'fluent:quiz-20-regular',
          affix: true,
        },
      },
    ],
    meta: { order: 2 },
  },
  {
    name: 'System-dev',
    path: '/system',
    component: Layout,
    children: [
      {
        path: '',
        component: () => import('@/views/system/index.vue'),
        name: 'System',
        meta: {
          title: 'System Default',
          icon: 'icon-park-outline:workbench',
          affix: true,
        },
      },
    ],
    meta: { order: 5 },
  },
  {
    name: 'Course-dev',
    path: '/course',
    component: Layout,
    children: [
      {
        path: '',
        component: () => import('@/views/course/index.vue'),
        name: 'Course',
        meta: {
          title: 'Course Default',
          icon: 'icon-park-outline:workbench',
          affix: true,
        },
      },
    ],
    meta: { order: 5 },
  },
    {
      name: 'Profile-dev',
      path: '/profile',
      component: Layout,
      isHidden: true,
      children: [
        {
          path: '',
          component: () => import('@/views/profile/index.vue'),
          name: 'Profile',
          meta: {
            title: 'Profile',
            icon: 'user',
            affix: true,
          },
        },
      ],
      meta: { order: 99 },
    },
]
