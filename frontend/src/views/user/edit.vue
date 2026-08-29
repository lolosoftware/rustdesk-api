<template>
  <div class="form-card">
    <el-form ref="root" label-width="120px" :model="form" :rules="rules">
      <el-form-item :label="T('Username')" prop="username">
        <el-input v-model="form.username"></el-input>
      </el-form-item>
      <el-form-item :label="T('Email')" prop="email">
        <el-input v-model="form.email"></el-input>
      </el-form-item>
      <el-form-item :label="T('Nickname')" prop="nickname">
        <el-input v-model="form.nickname"></el-input>
      </el-form-item>
      <el-form-item :label="T('Group')" prop="group_id">
        <el-select v-model="form.group_id">
          <el-option
              v-for="item in groupsList"
              :key="item.id"
              :label="item.name"
              :value="item.id"
          ></el-option>
        </el-select>
      </el-form-item>
      <el-form-item :label="T('IsAdmin')" prop="is_admin">
        <el-switch v-model="form.is_admin"
                   :active-value="true"
                   :inactive-value="false"
        ></el-switch>
      </el-form-item>
      <el-form-item :label="T('Status')" prop="status">
        <el-switch v-model="form.status"
                   :active-value="ENABLE_STATUS"
                   :inactive-value="DISABLE_STATUS"
        ></el-switch>
      </el-form-item>
      <el-form-item label="Activer OTP">
        <el-switch v-model="form.otp_enabled" @change="ensureOtpSecret"></el-switch>
      </el-form-item>
      <el-form-item v-if="form.otp_enabled" label="Secret OTP">
        <el-input v-model="form.otp_secret" readonly>
          <template #append>
            <el-button @click="copyOtpSecret">Copier</el-button>
          </template>
        </el-input>
        <div class="otp-help">Ajoutez ce secret dans votre application d’authentification.</div>
      </el-form-item>
      <el-form-item :label="T('Remark')" prop="remark">
          <el-input v-model="form.remark"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button @click="cancel">{{ T('Cancel') }}</el-button>
        <el-button @click="submit" type="primary">{{ T('Submit') }}</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
  import { useRoute } from 'vue-router'
  import { useGetDetail, useSubmit } from '@/views/user/composables/edit'
  import { ENABLE_STATUS, DISABLE_STATUS } from '@/utils/common_options'
  import { T } from '@/utils/i18n'
  import { ElMessage } from 'element-plus'

  const route = useRoute()
  const { form, item, getDetail, groupsList } = useGetDetail(route.params.id)

  const { root, rules, validate, submit, cancel } = useSubmit(form, route.params.id)

  const ensureOtpSecret = (enabled) => {
    if (!enabled || form.value.otp_secret) return
    const alphabet = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ234567'
    const bytes = new Uint8Array(20)
    crypto.getRandomValues(bytes)
    let secret = ''
    let buffer = 0
    let bits = 0
    for (const byte of bytes) {
      buffer = (buffer << 8) | byte
      bits += 8
      while (bits >= 5) {
        secret += alphabet[(buffer >>> (bits - 5)) & 31]
        bits -= 5
      }
    }
    if (bits > 0) secret += alphabet[(buffer << (5 - bits)) & 31]
    form.value.otp_secret = secret
  }

  const copyOtpSecret = async () => {
    await navigator.clipboard.writeText(form.value.otp_secret)
    ElMessage.success('Secret OTP copié')
  }

</script>

<style lang="scss" scoped>
.form-card {
}
.otp-help {
  color: #909399;
  font-size: 12px;
  line-height: 1.5;
}
</style>
