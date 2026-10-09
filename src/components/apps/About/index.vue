<script setup lang="ts">
import { NButton, NDivider, NGradientText, NTag, useMessage } from 'naive-ui'
import { onMounted, ref } from 'vue'
import { checkNewVersion, get, runningInfo } from '@/api/system/about'
import { t } from '@/locales'
import srcSvglogo from '@/assets/logo.svg'
import srcGithub from '@/assets/about_image/github.png'
import srcDocker from '@/assets/about_image/docker.png'

interface Version {
  versionName: string
  versionCode: number
}

interface RunningInfo {
  versionName: string
  versionCode: number
  goVersion: string
  os: string
  arch: string
  isDocker: boolean
  database: string
  runTime: string
  startTime: string
  userCount: number
  groupCount: number
  itemCount: number
  fileCount: number
  fileSize: number
}

const ms = useMessage()
const versionName = ref('')
const frontVersion = import.meta.env.VITE_APP_VERSION || 'unknown'

const running = ref<RunningInfo | null>(null)
const checkVersionLoading = ref(false)
const newVersion = ref('')

onMounted(() => {
  get<Version>().then((res) => {
    if (res.code === 0)
      versionName.value = res.data.versionName
  })
  runningInfo<RunningInfo>().then((res) => {
    if (res.code === 0)
      running.value = res.data
  }).catch(() => {})
})

async function handleCheckNewVersion() {
  checkVersionLoading.value = true
  try {
    const { code, data } = await checkNewVersion<{ hasNewVersion: boolean, newVersion: string, nowVersion: string }>(false)
    if (code === 0) {
      if (data?.hasNewVersion) {
        newVersion.value = data.newVersion
        ms.success(`发现新版本：${data.newVersion}`)
      }
      else {
        ms.success(t('apps.about.alreadyLatest'))
      }
    }
  }
  catch {
    ms.error(t('common.networkError'))
  }
  checkVersionLoading.value = false
}

function formatSize(size: number) {
  if (!size)
    return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0
  let v = size
  while (v >= 1024 && i < units.length - 1) {
    v /= 1024
    i++
  }
  return `${v.toFixed(i === 0 ? 0 : 1)} ${units[i]}`
}
</script>

<template>
  <div class="pt-5">
    <div class="flex flex-col items-center justify-center">
      <img :src="srcSvglogo" width="100" height="100" alt="">
      <div class="text-3xl font-semibold">
        {{ $t('common.appName') }}
      </div>
      <div class="text-xl">
        <NGradientText type="info">
          <a href="https://github.com/zhf883680/sun-panel/releases" class="font-semibold" :title="$t('apps.about.viewUpdateLog')" target="_blank">v{{ versionName }}</a>
        </NGradientText>
      </div>
      <div class="mt-2">
        <a href="https://github.com/zhf883680/sun-panel/releases" target="_blank" class="link">{{ $t('apps.about.checkUpdate') }}</a>
        <NButton class="ml-2" size="tiny" :loading="checkVersionLoading" @click="handleCheckNewVersion">
          {{ $t('apps.about.checkNewVersion') }}
        </NButton>
        <NTag v-if="newVersion" class="ml-2" size="tiny" type="success" :bordered="false">
          {{ newVersion }}
        </NTag>
      </div>
    </div>

    <NDivider style="margin:10px 0">
      •
    </NDivider>
    <div class="flex flex-col items-center justify-center text-base">
      <div class="text-xs opacity-70">
        {{ $t('apps.about.forkNote') }}
      </div>
      <div class="text-xs opacity-70 mt-[4px]">
        {{ $t('apps.about.originalRepo') }}
        <a href="https://github.com/hslr-s/sun-panel" target="_blank" class="link">github.com/hslr-s/sun-panel</a>
        ·
        {{ $t('apps.about.originalAuthor') }}
        <a href="https://github.com/hslr-s" target="_blank" class="link">红烧猎人</a>
      </div>

      <div class="mt-[8px]">
        {{ $t('apps.about.issue') }}<a href="https://github.com/zhf883680/sun-panel/issues" target="_blank" class="link">Github Issues</a>
      </div>
      <div>
        {{ $t('apps.about.discussions') }}<a href="https://github.com/zhf883680/sun-panel/discussions" target="_blank" class="link">Github Discussions</a>
      </div>

      <div class="flex mt-[10px] flex-wrap justify-center">
        <div class="flex items-center mx-[10px]">
          <img class="w-[20px] h-[20px] mr-[5px]" :src="srcGithub" alt="">
          <a href="https://github.com/zhf883680/sun-panel" target="_blank" class="link">Github</a>
        </div>
        <div class="flex items-center mx-[10px]">
          <img class="w-[20px] h-[20px] mr-[5px]" :src="srcDocker" alt="">
          <a href="https://hub.docker.com/r/zhf883680/sun-panel" target="_blank" class="link">Docker</a>
        </div>
      </div>

      <div class="mt-5">
        <NTag :bordered="false" size="small">
          {{ $t("apps.about.frontVersionText") }}: FV-{{ frontVersion }}
        </NTag>
      </div>

      <div v-if="running" class="mt-5 w-full max-w-[460px] text-left">
        <NDivider style="margin:10px 0">•</NDivider>
        <div class="text-center font-bold mb-2">{{ $t('apps.about.runningInfo') }}</div>
        <div class="grid grid-cols-2 gap-x-4 gap-y-1 text-xs">
          <div class="flex justify-between"><span class="opacity-60">{{ $t('apps.about.backendVersion') }}</span><span>{{ (running.versionName || '').trim() }}</span></div>
          <div class="flex justify-between"><span class="opacity-60">{{ $t('apps.about.runTime') }}</span><span>{{ running.runTime }}</span></div>
          <div class="flex justify-between"><span class="opacity-60">{{ $t('apps.about.goVersion') }}</span><span>{{ running.goVersion }}</span></div>
          <div class="flex justify-between"><span class="opacity-60">{{ $t('apps.about.startTime') }}</span><span>{{ running.startTime }}</span></div>
          <div class="flex justify-between"><span class="opacity-60">{{ $t('apps.about.osArch') }}</span><span>{{ running.os }}/{{ running.arch }}</span></div>
          <div class="flex justify-between"><span class="opacity-60">{{ $t('apps.about.database') }}</span><span>{{ running.database }}</span></div>
          <div class="flex justify-between"><span class="opacity-60">{{ $t('apps.about.userCount') }}</span><span>{{ running.userCount }}</span></div>
          <div class="flex justify-between"><span class="opacity-60">{{ $t('apps.about.groupCount') }}</span><span>{{ running.groupCount }}</span></div>
          <div class="flex justify-between"><span class="opacity-60">{{ $t('apps.about.itemCount') }}</span><span>{{ running.itemCount }}</span></div>
          <div class="flex justify-between"><span class="opacity-60">{{ $t('apps.about.fileCount') }}</span><span>{{ running.fileCount }}</span></div>
          <div class="flex justify-between"><span class="opacity-60">{{ $t('apps.about.fileSize') }}</span><span>{{ formatSize(running.fileSize) }}</span></div>
          <div class="flex justify-between"><span class="opacity-60">{{ $t('apps.about.isDocker') }}</span><span>{{ running.isDocker ? 'Yes' : 'No' }}</span></div>
        </div>
      </div>

    </div>
  </div>
</template>

<style>
.link{
    color:rgb(0, 89, 255)
}
</style>
