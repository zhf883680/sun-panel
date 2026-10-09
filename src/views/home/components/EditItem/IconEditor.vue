<script setup lang="ts">
import { NButton, NInput, NRadio, NSpace, NUpload } from 'naive-ui'
import type { UploadFileInfo } from 'naive-ui'
import { computed, defineProps, ref } from 'vue'
import { GallerySelector, ItemIcon } from '@/components/common'
import { useAuthStore } from '@/store'
import { apiRespErrMsg } from '@/utils/request/apiMessage'

const props = withDefaults(defineProps<{
  itemIcon: Panel.ItemIcon | null
  // 卡片背景色：预览时直接套用，避免与“卡片背景色”配置出现两个入口
  cardBackground?: string
}>(), {
  cardBackground: '#2a2a2a6b',
})
const emit = defineEmits<{
  (e: 'update:itemIcon', visible: Panel.ItemIcon): void // 定义修改父组件（prop内）的值的事件
}>()
const authStore = useAuthStore()
const showGallery = ref(false)

const initData: Panel.ItemIcon = {
  itemType: 2,
  backgroundColor: '#2a2a2a6b',
}

const itemIconInfo = computed({
  get() {
    const v = {
      ...initData,
      ...props.itemIcon,
      backgroundColor: props.itemIcon?.backgroundColor || initData.backgroundColor,
    }
    return v
  },
  set() {
    handleChange()
  },
})

function handleIconTypeRadioChange(type: number) {
  // checkedValueRef.value = type
  itemIconInfo.value.itemType = type
  handleChange()
}

function handleChange() {
  emit('update:itemIcon', itemIconInfo.value || null)
}

// 从图库选择图片
function handleGallerySelect(src: string) {
  itemIconInfo.value.src = src
  handleChange()
}

const handleUploadFinish = ({
  file,
  event,
}: {
  file: UploadFileInfo
  event?: ProgressEvent
}) => {
  const res = JSON.parse((event?.target as XMLHttpRequest).response)
  if (res.code === 0) {
    const imageUrl = res.data.imageUrl
    itemIconInfo.value.src = imageUrl
    emit('update:itemIcon', itemIconInfo.value || null)
  }
  else {
    apiRespErrMsg(res)
    // ms.error(`${t('common.uploadFail')}:${res.msg}`)
  }

  return file
}
</script>

<template>
  <div>
    <div class="mb-[10px]">
      <NRadio
        :checked="itemIconInfo.itemType === 1 "
        :value="1"
        name="iconType"
        @change="handleIconTypeRadioChange(1)"
      >
        {{ $t('common.text') }}
      </NRadio>

      <NRadio
        :checked="itemIconInfo.itemType === 2"
        :value="2"
        name="iconType"
        @change="handleIconTypeRadioChange(2)"
      >
        {{ $t('common.image') }}
      </NRadio>

      <NRadio
        :checked="itemIconInfo.itemType === 3"
        :value="3"
        name="iconType"
        @change="handleIconTypeRadioChange(3)"
      >
        {{ $t('iconItem.onlineIcon') }}
      </NRadio>
    </div>

    <div class=" h-[100px]">
      <div class="flex">
        <div>
          <div class="border rounded-2xl bg-slate-200 overflow-hidden rounded-2xl transparent-grid">
            <ItemIcon :item-icon="itemIconInfo" :force-background="props.cardBackground" />
          </div>
        </div>
        <!-- 文字 -->
        <div class="ml-[20px]">
          <!-- <NImage :src="model.icon" preview-disabled /> -->
          <div v-if="itemIconInfo.itemType === 1">
            <NInput v-model:value="itemIconInfo.text" class="mb-[5px]" size="small" type="text" @input="handleChange" />
          </div>

          <div v-if="itemIconInfo.itemType === 3">
            <div>
              <NInput v-model:value="itemIconInfo.text" class="mb-[5px]" size="small" type="text" :placeholder="$t('iconItem.inputIconName')" @input="handleChange" />

              <NButton quaternary type="info">
                <a target="_blank" href="https://icon-sets.iconify.design/">{{ $t('iconItem.onlineIconLibrary') }}</a>
              </NButton>
            </div>
          </div>

          <!-- 图片 -->
          <div v-if="itemIconInfo.itemType === 2">
            <NInput v-model:value="itemIconInfo.src" class="mb-[5px] w-full" size="small" type="text" :placeholder="$t('iconItem.inputIconUrlOrUpload')" @input="handleChange" />
            <NSpace :size="8">
              <NButton size="small" type="primary" @click="showGallery = true">
                {{ $t('iconItem.gallery') }}
              </NButton>
              <NUpload
                action="/api/file/uploadImg"
                :show-file-list="false"
                name="imgfile"
                :data="{ category: 'icon' }"
                :headers="{
                  token: authStore.token as string,
                }"
                @finish="handleUploadFinish"
              >
                <NButton size="small">
                  {{ $t('iconItem.selectUpload') }}
                </NButton>
              </NUpload>
            </NSpace>
          </div>
        </div>
      </div>
    </div>

    <GallerySelector v-model:show="showGallery" default-type="icon" @select="handleGallerySelect" />
  </div>
</template>

<style scoped>
.transparent-grid {
    background-image: linear-gradient(45deg, #e6e4e4 25%, transparent 25%, transparent 75%, #e6e4e4 75%),
                      linear-gradient(45deg, #e6e4e4 25%, transparent 25%, transparent 75%, #e6e4e4 75%);
    background-size: 16px 16px;
    background-position: 0 0, 8px 8px;
}
</style>
