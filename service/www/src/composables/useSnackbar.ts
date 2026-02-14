import { reactive } from 'vue'

const snackbar = reactive({
  show: false,
  text: '',
  color: 'success',
})

export function useSnackbar () {
  function showNotice (text: string, color = 'success') {
    snackbar.text = text
    snackbar.color = color
    snackbar.show = true
  }

  return { snackbar, showNotice }
}
