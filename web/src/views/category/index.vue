<script setup lang="ts">
import { NButton } from 'naive-ui'
// hooks
import { useModalStore } from '@/store/modal'
import { useCategoriesQuery } from '@/data/category'
// components
import TheIcon from '@/components/icon/TheIcon.vue'
import CommonPage from '@/components/page/CommonPage.vue'
import CategoryList from '@/components/category/CategoryList.vue'
import CategoryModal from '@/components/category/CategoryModal.vue'

// query
const { categories, loading } = useCategoriesQuery({})
// store hooks
const modal = useModalStore()

function openCreateModal() {
  modal.open(CategoryModal, {
    title: 'Create Category',
  })
}
</script>

<template>
  <CommonPage show-footer title="Category List">
    <template #action>
      <div>
        <NButton class="float-right mr-15" type="primary" @click="openCreateModal">
          <TheIcon icon="material-symbols:add" :size="18" class="mr-5" />Create new category
        </NButton>
      </div>
    </template>
    <CategoryList :loading="loading" :table-data="categories" />
  </CommonPage>
</template>
