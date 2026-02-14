<script setup lang="ts">
import { ref, watch } from 'vue'
import { NForm, NFormItem, NInput, NButton } from 'naive-ui'
// hooks
import { useModalStore } from '@/store/modal'
import { useCreateCategoryMutation, useUpdateCategoryMutation } from '@/data/category'
// types
import type { Category } from '@/types'

const props = defineProps<{
  category?: Category | null
}>()

const modal = useModalStore()

// mutations
const { mutateAsync: createCategory, isPending: creating } = useCreateCategoryMutation()
const { mutateAsync: updateCategory, isPending: updating } = useUpdateCategoryMutation()

const modalFormRef = ref()
const modalForm = ref({
  name: '',
})

watch(
  () => props.category,
  (category) => {
    if (!category) {
      // create mode
      modalForm.value = {
        name: '',
      }
      return
    }

    // edit mode
    modalForm.value = {
      name: category.name,
    }
  },
  { immediate: true },
)

const validateAddCourse = {
  name: [{ required: true, message: 'Name is required', trigger: ['blur'] }],
  serialNumber: [{ required: true, message: 'Serial Number is required', trigger: ['blur'] }],
}

async function handleSave() {
  modalFormRef.value?.validate(async (errors: Error) => {
    if (errors) return
    if (props.category) {
      await updateCategory({
        id: props.category.id,
        name: modalForm.value.name,
      })
    } else {
      await createCategory({
        name: modalForm.value.name,
        description: null,
      })
    }
    modal.close()
  })
}
</script>

<template>
  <div>
    <!-- FORM -->
    <NForm
      ref="modalFormRef"
      label-placement="left"
      label-align="left"
      :label-width="80"
      :model="modalForm"
      :rules="validateAddCourse"
    >
      <NFormItem label="Name" path="name">
        <NInput v-model:value="modalForm.name" />
      </NFormItem>

      <!-- <NFormItem label="Serial No." path="serialNumber">
        <NInput v-model:value="modalForm.serialNumber" />
      </NFormItem> -->
    </NForm>

    <div flex justify-end>
      <NButton @click="modal.close">Cancel</NButton>
      <NButton type="primary" class="ml-16" :loading="creating || updating" @click="handleSave">
        Save
      </NButton>
    </div>
  </div>
</template>
