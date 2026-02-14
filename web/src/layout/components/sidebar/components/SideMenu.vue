<template>
  <n-menu
    ref="menu"
    class="side-menu"
    accordion
    :indent="18"
    :collapsed-icon-size="22"
    :collapsed-width="64"
    :options="menuOptions"
    :value="activeKey"
    @update:value="handleMenuSelect"
  />
</template>

<script setup lang="ts">
import { renderCustomIcon, renderIcon } from '@/utils'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const router = useRouter()
const curRoute = useRoute()

const activeKey = computed(() => curRoute.meta?.activeMenu || curRoute.name)

const menuOptions = [
  {
    label: 'Workbench',
    key: 'workbench',
    icon: renderIcon('icon-park-outline:workbench', { size: 18 }),
    path: '/workbench',
  },
  {
    label: 'Category',
    key: 'category',
    icon: renderIcon('carbon:category', { size: 18 }),
    path: '/categories',
  },
  {
    label: 'Asset',
    key: 'aseet',
    icon: renderIcon('qlementine-icons:items-tree-16', { size: 18 }),
    path: '/assets',
  },
  // {
  //   label: 'Course',
  //   key: 'course',
  //   icon: renderIcon('material-symbols:book-outline', { size: 18 }),
  //   path: '/course',
  // },
  // {
  //   label: 'Quizzes',
  //   key: 'quizzes',
  //   icon: renderIcon('fluent:quiz-20-regular', { size: 20 }),
  //   path: '/quizzes',
  // },
  {
    label: 'System',
    key: 'system',
    icon: renderIcon('mdi-account-off', { size: 18 }),
    path: '/system',
  },
  // {
  //   label: 'Profile',
  //   key: 'profile',
  //   icon: renderIcon('mdi-account-off', { size: 18 }),
  //   path: '/profile',
  // },
  // {
  //   label: "Settings",
  //   key: "settings",
  //   children: [
  //     { label: "Profile", key: "profile", path: '/profile' },
  //     { label: "Security", key: "security", path: '/security'  },
  //   ],
  // },
]

function handleMenuSelect(key: number, item: any) {
  router.push(item.path)
}
</script>

<style lang="scss">
.side-menu:not(.n-menu--collapsed) {
  .n-menu-item-content {
    &::before {
      left: 5px;
      right: 5px;
    }
    &.n-menu-item-content--selected,
    &:hover {
      &::before {
        border-left: 4px solid var(--primary-color);
      }
    }
  }
}
</style>
