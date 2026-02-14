<script setup lang="ts">
import { h } from 'vue'
import { NButton, NPopconfirm } from 'naive-ui'
// hooks
import { useDeleteCategoryMutation } from '@/data/category'
// utils
import { renderIcon } from '@/utils'
// types
import type { Category } from '@/types'
// hooks
import { useModalStore } from '@/store/modal'
// components
import CategoryModal from './CategoryModal.vue'

// Define props
const props = defineProps<{
  loading: boolean
  tableData: Category[]
}>()
const modal = useModalStore()

// mutation
const { mutateAsync: deleteCategory } = useDeleteCategoryMutation()

function onEdit(category: any) {
  modal.open(CategoryModal, {
    title: 'Edit Category',
    props: {
      category,
    },
  })
}

async function deleteRow(row: Category) {
  await deleteCategory({ id: row.id })
}

const columns = [
  {
    title: 'Name',
    key: 'name',
    width: 200,
  },
  {
    title: 'Actions',
    key: 'actions',
    width: 160,
    render(row: any) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            class: 'mr-5',
            onClick: () => onEdit(row),
          },
          { default: () => 'Edit', icon: renderIcon('material-symbols:edit', { size: 16 }) },
        ),
        h(
          NPopconfirm,
          {
            onPositiveClick: () => deleteRow(row),
          },
          {
            trigger: () =>
              h(
                NButton,
                {
                  size: 'small',
                  type: 'error',
                },
                {
                  default: () => 'Delete',
                  icon: renderIcon('material-symbols:delete-outline', { size: 16 }),
                },
              ),
            default: () => 'Are you sure?',
          },
        ),
      ]
    },
  },
]
</script>

<template>
  <n-data-table :loading="loading" :columns="columns" :data="tableData" />
</template>
