<template>
  <div style="margin-bottom: 50px;font-weight: bold">灵敢足迹导出csv工具</div>
  <div style="display: flex;width: 500px;justify-content: center;">
    <el-upload
        ref="upload"
        class="upload-demo"
        :action="uploadUrl"
        :limit="1"
        :on-exceed="handleExceed"
        :auto-upload="false"
        :on-change="handleUpload"
    >
      <template #trigger>
        <el-button type="primary" style="margin-right: 20px">选择文件</el-button>
      </template>
      <template #default>
        <el-button type="success" @click="submitUpload">
          上传
        </el-button>
      </template>
      <template #tip>
        <div class="el-upload__tip text-red">
          上传新文件将会替换旧文件
        </div>
      </template>
    </el-upload>
  </div>
</template>

<script setup lang="ts">
import {ref} from 'vue'
import {genFileId, UploadFile} from 'element-plus'
import type {UploadInstance, UploadProps, UploadRawFile} from 'element-plus'
import {ResultRes} from '@/utils/beans.ts';
import {ElNotification as nofity} from 'element-plus'
import {useRouter} from "vue-router";
import apiUrl from "@/config/config.ts";
import {getApiBaseUrl} from "@/axios";

const upload = ref<UploadInstance>()
const $router = useRouter()

const uploadUrl = getApiBaseUrl(9091) + '/upload'

const handleExceed: UploadProps['onExceed'] = (files) => {
  upload.value!.clearFiles()
  const file = files[0] as UploadRawFile
  file.uid = genFileId()
  upload.value!.handleStart(file)
}

const submitUpload = () => {
  upload.value!.submit()
}

const handleUpload = (uploadFile: UploadFile) => {
  if (uploadFile.status === 'success') {
    const res: ResultRes = uploadFile.response
    if (res?.success) {
      sessionStorage.setItem('session', res.data?.session)
      nofity.success({
        message: '上传成功',
      })
      $router.push({ name:'Home' })
    } else {
      nofity.warning({
        message: res?.msg,
      })
    }
  }
  if (uploadFile.status === 'fail') {
    nofity.error({
      message: '上传失败，请检查后台',
    })
  }
}
</script>
