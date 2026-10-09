<script setup lang="ts">
import { NAlert, NButton, NCard, NFormItem, NImage, NInput, NPopconfirm, NSpace, NSwitch, NUpload, useMessage } from 'naive-ui'
import { onMounted, ref } from 'vue'
import {
  getCustomJsAndCssCode,
  getCustomStylePath,
  getSiteStting,
  saveCustomJsAndCssCode,
  saveSiteStting,
} from '@/api/panel/globalSetting'
import type { SiteSetting } from '@/api/panel/globalSetting'
import { GallerySelector } from '@/components/common'
import { useAppStore, useAuthStore } from '@/store'
import { defaultSiteSetting } from '@/store/modules/app/helper'
import { t } from '@/locales'

const ms = useMessage()
const appStore = useAppStore()
const authStore = useAuthStore()

const siteSetting = ref<SiteSetting>({ ...defaultSiteSetting() })
const saveLoading = ref(false)

const customPath = ref('')
const customCode = ref({ jsContent: '', cssContent: '' })
const customSaveLoading = ref(false)

// 站点图标 / 登录页背景 的图库选择
const showGallery = ref(false)
const galleryTarget = ref<'favicon' | 'loginBg'>('favicon')

async function getSiteSetting() {
  try {
    const { code, data } = await getSiteStting<SiteSetting>()
    if (code === 0 && data)
      siteSetting.value = { ...defaultSiteSetting(), ...data }
  }
  catch {
    ms.error(t('common.failed'))
  }
}

async function handleSaveSiteSetting() {
  saveLoading.value = true
  try {
    const { code, msg } = await saveSiteStting(siteSetting.value)
    if (code === 0) {
      appStore.setSiteSetting(siteSetting.value)
      ms.success(t('common.saveSuccess'))
    }
    else {
      ms.error(`${t('common.saveFail')}:${msg}`)
    }
  }
  catch {
    ms.error(t('common.saveFail'))
  }
  saveLoading.value = false
}

async function getCustomCode() {
  try {
    const { data: pathData } = await getCustomStylePath<{ customPath: string }>()
    customPath.value = pathData?.customPath || ''
  }
  catch { /* ignore */ }
  try {
    const { data } = await getCustomJsAndCssCode<{ jsContent: string, cssContent: string }>()
    customCode.value = { jsContent: data?.jsContent || '', cssContent: data?.cssContent || '' }
  }
  catch { /* ignore */ }
}

async function handleSaveCustomCode() {
  customSaveLoading.value = true
  try {
    const { code, msg } = await saveCustomJsAndCssCode(customCode.value)
    if (code === 0)
      ms.success(t('common.saveSuccess'))
    else
      ms.error(`${t('common.saveFail')}:${msg}`)
  }
  catch {
    ms.error(t('common.saveFail'))
  }
  customSaveLoading.value = false
}

function handleResetCustomCode() {
  customCode.value = { jsContent: '', cssContent: '' }
}

function handleSelectFromGallery() {
  showGallery.value = true
}

function handleGallerySelect(src: string) {
  if (galleryTarget.value === 'favicon')
    siteSetting.value.favicoUrl = src
  else
    siteSetting.value.loginBackgroundUrl = src
}

async function handleUploadFavicon({ file, event }: any) {
  const res = JSON.parse(event?.target?.response || '{}')
  if (res.code === 0) {
    siteSetting.value.favicoUrl = res.data.imageUrl
    ms.success(t('common.uploadSuccess'))
  }
  else {
    ms.error(`${t('common.uploadFail')}:${res.msg}`)
  }
  return file
}

onMounted(() => {
  getSiteSetting()
  getCustomCode()
})
</script>

<template>
  <div class="bg-slate-200 dark:bg-zinc-900 rounded-[10px] p-[8px] overflow-auto h-full">
    <!-- 站点设置 -->
    <NCard :bordered="false" size="small" style="border-radius: 10px;">
      <template #header>
        {{ $t('apps.globalSetting.personalise') }}
      </template>

      <NFormItem :label="$t('apps.globalSetting.title')">
        <NInput v-model:value="siteSetting.title" type="text" show-count :maxlength="20" :placeholder="$t('common.inputPlaceholder')" />
      </NFormItem>

      <NFormItem :label="$t('apps.globalSetting.favicon')">
        <div class="w-full">
          <div class="flex items-center gap-2">
            <NImage
              v-if="siteSetting.favicoUrl"
              width="40" height="40" style="border:1px solid #e9e9e9;border-radius:8px;"
              :src="siteSetting.favicoUrl" preview-disabled
            />
            <NSpace :size="8">
              <NButton size="small" @click="galleryTarget = 'favicon'; handleSelectFromGallery()">
                {{ $t('common.selectFromGallery') }}
              </NButton>
              <NUpload
                action="/api/panel/globalSetting/uploadFaviconImage"
                :show-file-list="false" name="imgfile"
                :headers="{ token: authStore.token as string }"
                accept=".ico,.png,.svg"
                @finish="handleUploadFavicon"
              >
                <NButton size="small">{{ $t('apps.globalSetting.uploadIcoImage') }}</NButton>
              </NUpload>
              <NButton size="small" :disabled="!siteSetting.favicoUrl" @click="siteSetting.favicoUrl = ''">
                {{ $t('apps.globalSetting.clearIcoImage') }}
              </NButton>
            </NSpace>
          </div>
          <NInput v-model:value="siteSetting.favicoUrl" class="mt-2" size="small" type="text" :placeholder="$t('apps.globalSetting.customImageURL')" />
        </div>
      </NFormItem>

      <NFormItem :label="$t('apps.globalSetting.loginBackgroundImage')">
        <div class="w-full">
          <NSpace :size="8">
            <NButton size="small" @click="galleryTarget = 'loginBg'; handleSelectFromGallery()">
              {{ $t('common.selectFromGallery') }}
            </NButton>
            <NButton v-if="siteSetting.loginBackgroundUrl" size="small" @click="siteSetting.loginBackgroundUrl = ''">
              {{ $t('common.reset') }}
            </NButton>
          </NSpace>
          <NInput v-model:value="siteSetting.loginBackgroundUrl" class="mt-2" size="small" type="text" :placeholder="$t('apps.globalSetting.customImageURL')" />
        </div>
      </NFormItem>

      <NFormItem :label="$t('apps.globalSetting.isDisableItemCardCache')">
        <div class="w-full">
          <NSwitch v-model:value="siteSetting.isDisableItemCardCache" />
          <div class="text-xs text-slate-400 mt-[2px]">
            {{ $t('apps.globalSetting.isDisableItemCardCacheTip') }}
          </div>
        </div>
      </NFormItem>

      <div class="flex justify-end">
        <NButton type="primary" :loading="saveLoading" @click="handleSaveSiteSetting">
          {{ $t('common.save') }}
        </NButton>
      </div>
    </NCard>

    <!-- 自定义JS和CSS -->
    <NCard :bordered="false" size="small" class="mt-3" style="border-radius: 10px;">
      <template #header>
        {{ $t('apps.globalSetting.customJSAndCss') }}
      </template>

      <NAlert type="error" :bordered="false" style="border-radius: 5px;">
        {{ $t('apps.globalSetting.customDangerAlert') }}
      </NAlert>

      <div class="text-xs text-slate-400 my-[6px]">
        {{ $t('apps.globalSetting.customJSAndCssAlert', { path: customPath }) }}
      </div>

      <div class="text-slate-500 mb-[5px] font-bold">
        JS ({{ customPath }}/index.js)
      </div>
      <NInput v-model:value="customCode.jsContent" type="textarea" :rows="6" clearable />

      <div class="text-slate-500 mb-[5px] mt-[10px] font-bold">
        CSS ({{ customPath }}/index.css)
      </div>
      <NInput v-model:value="customCode.cssContent" type="textarea" :rows="6" clearable />

      <div class="flex justify-end mt-3 gap-2">
        <NPopconfirm @positive-click="handleResetCustomCode">
          <template #trigger>
            <NButton size="small" type="error">{{ $t('common.reset') }}</NButton>
          </template>
          {{ $t('apps.globalSetting.resetWarnText') }}
        </NPopconfirm>
        <NButton size="small" type="success" :loading="customSaveLoading" @click="handleSaveCustomCode">
          {{ $t('common.save') }}
        </NButton>
      </div>
    </NCard>

    <GallerySelector v-model:show="showGallery" default-type="icon" @select="handleGallerySelect" />
  </div>
</template>
