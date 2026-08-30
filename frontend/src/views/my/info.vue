<template>
  <div>
    <el-card :title="T('Userinfo')" shadow="hover">
      <el-form class="info-form" ref="form" label-width="120px" label-suffix="：">
        <el-form-item :label="T('Username')">
          <div>{{ userStore.username }}</div>
        </el-form-item>
        <el-form-item :label="T('Email')">
          <div>{{ userStore.email }}</div>
        </el-form-item>
        <el-form-item :label="T('Password')" prop="password">
          <el-button type="danger" @click="showChangePwd">{{ T('ChangePassword') }}</el-button>
        </el-form-item>
        <el-form-item :label="T('TwoFactorAuthentication')">
          <el-tag v-if="otpEnabled" type="success">{{ T('Enabled') }}</el-tag>
          <el-tag v-else type="info">{{ T('Disabled') }}</el-tag>
          <el-button v-if="!otpEnabled" type="primary" style="margin-left: 12px" @click="beginOtpSetup">
            {{ T('Configure') }}
          </el-button>
          <el-button v-else type="danger" style="margin-left: 12px" @click="disableOtp">
            {{ T('Disable') }}
          </el-button>
        </el-form-item>
        <el-form-item label="OIDC">
          <el-table :data="oidcData" border fit>
            <el-table-column :label="T('IdP')" prop="op" align="center"></el-table-column>
            <el-table-column :label="T('Status')" prop="status" align="center">
              <template #default="{ row }">
                <el-tag v-if="row.status === 1" type="success">{{ T('HasBind') }}</el-tag>
                <el-tag v-else type="danger">{{ T('NoBind') }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column :label="T('Actions')" align="center" width="200">
              <template #default="{ row }">
                <el-button v-if="row.status === 1" type="danger" size="small" @click="toUnBind(row)">{{ T('UnBind') }}</el-button>
                <el-button v-else type="success" size="small" @click="toBind(row)">{{ T('ToBind') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-form-item>
      </el-form>
    </el-card>
    <el-card shadow="hover" style="margin-top: 20px">
      <div v-html="html"></div>
    </el-card>
    <changePwdDialog v-model:visible="changePwdVisible"></changePwdDialog>
    <el-dialog v-model="otpDialogVisible" :title="T('ConfigureTwoFactorAuthentication')" width="520px" :close-on-click-modal="false">
      <div class="otp-enrollment">
        <p>{{ T('ScanAuthenticatorQR') }}</p>
        <img :src="otpEnrollment.qr_code" :alt="T('AuthenticatorQRCode')" class="otp-qr"/>
        <p>{{ T('ManualAuthenticatorKey') }}</p>
        <el-input :model-value="otpEnrollment.secret" readonly>
          <template #append>
            <el-button @click="copyOtpSecret">{{ T('Copy') }}</el-button>
          </template>
        </el-input>
        <p>{{ T('EnterAuthenticatorCodeToConfirm') }}</p>
        <el-input v-model="otpCode" maxlength="6" inputmode="numeric" autocomplete="one-time-code"
                  @keyup.enter="confirmOtp"/>
      </div>
      <template #footer>
        <el-button @click="otpDialogVisible = false">{{ T('Cancel') }}</el-button>
        <el-button type="primary" :disabled="!otpCodeValid" @click="confirmOtp">
          {{ T('Confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import changePwdDialog from '@/components/changePwdDialog.vue'
  import { computed, ref } from 'vue'
  import { useUserStore } from '@/store/user'
  import { useAppStore } from '@/store/app'
  import { bind, unbind } from '@/api/oauth'
  import { myOauth, otpStatus, otpSetup, otpConfirm, otpDisable } from '@/api/user'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { T } from '@/utils/i18n'
  import { marked } from 'marked'

  const appStore = useAppStore()
  const userStore = useUserStore()
  const changePwdVisible = ref(false)
  const otpEnabled = ref(false)
  const otpDialogVisible = ref(false)
  const otpCode = ref('')
  const otpCodeValid = computed(() => /^[0-9]{6}$/.test(otpCode.value))
  const otpEnrollment = ref({ secret: '', qr_code: '', provisioning_uri: '' })
  const showChangePwd = () => {
    changePwdVisible.value = true
  }
  const oidcData = ref([])
  const getMyOauth = async () => {
    const res = await myOauth().catch(_ => false)
    if (res) {
      oidcData.value = res.data
    }

  }
  getMyOauth()
  const loadOtpStatus = async () => {
    const res = await otpStatus().catch(() => false)
    if (res) otpEnabled.value = Boolean(res.data.enabled)
  }
  loadOtpStatus()

  const beginOtpSetup = async () => {
    const res = await otpSetup().catch(() => false)
    if (!res) return
    otpEnrollment.value = res.data
    otpCode.value = ''
    otpDialogVisible.value = true
  }

  const copyOtpSecret = async () => {
    await navigator.clipboard.writeText(otpEnrollment.value.secret)
    ElMessage.success(T('Copied'))
  }

  const confirmOtp = async () => {
    if (!/^[0-9]{6}$/.test(otpCode.value)) return
    const res = await otpConfirm(otpCode.value).catch(() => false)
    if (!res) return
    otpDialogVisible.value = false
    otpEnabled.value = true
    userStore.otp_enabled = true
    ElMessage.success(T('TwoFactorAuthenticationEnabled'))
  }

  const disableOtp = async () => {
    const result = await ElMessageBox.prompt(
      T('EnterAuthenticatorCodeToDisable'),
      T('DisableTwoFactorAuthentication'),
      {
        inputPattern: /^[0-9]{6}$/,
        inputErrorMessage: T('AuthenticationCodeMustHaveSixDigits'),
        confirmButtonText: T('Confirm'),
        cancelButtonText: T('Cancel'),
      },
    ).catch(() => false)
    if (!result) return
    const res = await otpDisable(result.value).catch(() => false)
    if (!res) return
    otpEnabled.value = false
    userStore.otp_enabled = false
    ElMessage.success(T('TwoFactorAuthenticationDisabled'))
  }
  const toBind = async (row) => {
    const res = await bind({ op: row.op }).catch(_ => false)
    if (res) {
      const { code, url } = res.data
      window.open(url)
    }
  }
  const toUnBind = async (row) => {
    const cf = await ElMessageBox.confirm(T('Confirm?', { param: T('UnBind') }), {
      confirmButtonText: T('Confirm'),
      cancelButtonText: T('Cancel'),
      type: 'warning',
    }).catch(_ => false)
    if (!cf) {
      return false
    }
    const res = await unbind({ op: row.op }).catch(_ => false)
    if (res) {
      getMyOauth()
    }

  }

  const html = computed(_ => marked(appStore.setting.hello||''))

</script>

<style scoped lang="scss">
.info-form {
  width: 600px;
  margin: 0 auto;

}
.otp-enrollment {
  text-align: center;
}
.otp-qr {
  width: 256px;
  height: 256px;
  image-rendering: pixelated;
}
</style>
