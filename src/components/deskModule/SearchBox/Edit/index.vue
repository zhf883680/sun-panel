<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { NAvatar, NButton, NButtonGroup, NForm, NFormItem, NInput, NModal, useDialog, useMessage } from 'naive-ui'
import { SvgIcon } from '@/components/common'
import { t } from '@/locales'

interface Props {
  visible: boolean
  engineList: DeskModule.SearchBox.SearchEngine[]
  defaultList: DeskModule.SearchBox.SearchEngine[]
}

const props = withDefaults(defineProps<Props>(), {
  defaultList: () => [],
})
const emit = defineEmits<{
  (e: 'update:visible', visible: boolean): void
  (e: 'confirm', engineList: DeskModule.SearchBox.SearchEngine[]): void
}>()

const ms = useMessage()
const dialog = useDialog()

const show = computed({
  get: () => props.visible,
  set: (visible: boolean) => emit('update:visible', visible),
})

const list = ref<DeskModule.SearchBox.SearchEngine[]>([])
// null：列表视图；-1：新增；>=0：编辑
const editing = ref<number | null>(null)
const form = ref<DeskModule.SearchBox.SearchEngine>({ iconSrc: '', title: '', url: '' })
const formRef = ref<any>(null)

const rules = {
  title: {
    required: true,
    trigger: 'blur',
    message: t('deskModule.searchBox.engineTitleRequired'),
  },
  url: [
    {
      required: true,
      trigger: 'blur',
      message: t('deskModule.searchBox.engineUrlRequired'),
    },
    {
      trigger: 'blur',
      validator(rule: any, value: string) {
        if (value && !/^https?:\/\//.test(value))
          return new Error(t('deskModule.searchBox.engineUrlInvalid'))
        return true
      },
    },
  ],
}

const urlNoPlaceholder = computed(() => form.value.url !== '' && !form.value.url.includes('%s'))

watch(() => props.visible, (visible) => {
  if (visible) {
    list.value = JSON.parse(JSON.stringify(props.engineList || []))
    editing.value = null
  }
})

function handleAdd() {
  form.value = { iconSrc: '', title: '', url: '' }
  editing.value = -1
}

function handleEdit(index: number) {
  form.value = { ...list.value[index] }
  editing.value = index
}

function handleDelete(index: number) {
  dialog.warning({
    title: t('common.warning'),
    content: t('deskModule.searchBox.engineDeleteConfirm'),
    positiveText: t('common.confirm'),
    negativeText: t('common.cancel'),
    onPositiveClick: () => {
      list.value.splice(index, 1)
    },
  })
}

function handleMove(index: number, offset: number) {
  const target = index + offset
  if (target < 0 || target >= list.value.length)
    return
  const [item] = list.value.splice(index, 1)
  list.value.splice(target, 0, item)
}

// 恢复默认引擎
function handleReset() {
  list.value = JSON.parse(JSON.stringify(props.defaultList || []))
  ms.success(t('common.success'))
}

function handleFormConfirm() {
  formRef.value?.validate((errors: any) => {
    if (errors)
      return
    if (editing.value === -1)
      list.value.push({ ...form.value })
    else if (editing.value !== null)
      list.value[editing.value] = { ...form.value }
    editing.value = null
  })
}

function handleSave() {
  emit('confirm', JSON.parse(JSON.stringify(list.value)))
}
</script>

<template>
  <NModal v-model:show="show" preset="card" size="small" style="width: 620px;border-radius: 1rem;" :title="t('deskModule.searchBox.customEngine')">
    <!-- 列表视图 -->
    <div v-if="editing === null" class="max-h-[380px] overflow-auto p-[5px]">
      <div v-if="list.length === 0" class="text-center text-slate-400 py-[20px]">
        {{ t('deskModule.searchBox.engineEmpty') }}
      </div>

      <div
        v-for="engine, index in list" :key="index"
        class="flex items-center rounded-xl p-[8px] mb-[5px] bg-slate-100 dark:bg-zinc-800"
      >
        <div class="w-[36px] h-[36px] mr-[10px] flex items-center justify-center rounded-lg bg-white overflow-hidden">
          <NAvatar v-if="engine.iconSrc" :src="engine.iconSrc" style="background-color: transparent;" :size="24" />
          <span v-else class="text-xs text-slate-400">{{ (engine.title || '?').substring(0, 1) }}</span>
        </div>
        <div class="flex-1 overflow-hidden">
          <div class="font-bold truncate">
            {{ engine.title }}
          </div>
          <div class="text-xs text-slate-400 truncate">
            {{ engine.url }}
          </div>
        </div>
        <NButtonGroup>
          <NButton size="tiny" :disabled="index === 0" @click="handleMove(index, -1)">
            ↑
          </NButton>
          <NButton size="tiny" :disabled="index === list.length - 1" @click="handleMove(index, 1)">
            ↓
          </NButton>
          <NButton size="tiny" @click="handleEdit(index)">
            <template #icon>
              <SvgIcon icon="basil-edit-solid" />
            </template>
          </NButton>
          <NButton size="tiny" type="error" @click="handleDelete(index)">
            <template #icon>
              <SvgIcon icon="material-symbols-delete" />
            </template>
          </NButton>
        </NButtonGroup>
      </div>

      <NButton class="mt-[10px]" dashed block @click="handleAdd">
        <template #icon>
          <SvgIcon icon="typcn-plus" />
        </template>
        {{ t('deskModule.searchBox.addEngine') }}
      </NButton>
    </div>

    <!-- 编辑视图 -->
    <div v-else>
      <NForm ref="formRef" :model="form" :rules="rules">
        <NFormItem path="title" :label="t('deskModule.searchBox.engineTitle')">
          <NInput v-model:value="form.title" type="text" :maxlength="20" show-count />
        </NFormItem>
        <NFormItem path="url" :label="t('deskModule.searchBox.engineUrl')">
          <NInput v-model:value="form.url" type="text" :maxlength="200" placeholder="https://www.baidu.com/s?wd=%s" />
        </NFormItem>
        <div class="text-xs text-slate-400 mb-[10px]">
          {{ t('deskModule.searchBox.engineUrlTip') }}
        </div>
        <div v-if="urlNoPlaceholder" class="text-xs text-orange-400 mb-[10px]">
          {{ t('deskModule.searchBox.engineUrlNoPlaceholder') }}
        </div>
        <NFormItem path="iconSrc" :label="t('deskModule.searchBox.engineIcon')">
          <NInput v-model:value="form.iconSrc" type="text" :maxlength="500" placeholder="http(s)://" />
        </NFormItem>
      </NForm>

      <div class="flex justify-end">
        <NButton class="mr-[10px]" @click="editing = null">
          {{ t('common.cancel') }}
        </NButton>
        <NButton type="primary" @click="handleFormConfirm">
          {{ t('common.confirm') }}
        </NButton>
      </div>
    </div>

    <template v-if="editing === null" #footer>
      <div class="flex justify-between">
        <NButton @click="handleReset">
          {{ t('deskModule.searchBox.engineReset') }}
        </NButton>
        <div>
          <NButton class="mr-[10px]" @click="show = false">
            {{ t('common.cancel') }}
          </NButton>
          <NButton type="success" @click="handleSave">
            {{ t('common.save') }}
          </NButton>
        </div>
      </div>
    </template>
  </NModal>
</template>
