<script>
export default {
  data(){
    return { errormsg:null, msgs:[], body:'', loading:false, timer:null }
  },
  methods:{
    async load(){
      try{
        const id = localStorage.getItem('token')
        const peer = encodeURIComponent(this.$route.params.peer)
        const res = await this.$axios.get(`/users/${id}/chats/${peer}/messages`)
        this.msgs = res.data || []
      }catch(e){ this.errormsg = e.toString() }
    },
    async send(){
      if(!this.body.trim()) return
      try{
        const id = localStorage.getItem('token')
        const peer = encodeURIComponent(this.$route.params.peer)
        await this.$axios.post(`/users/${id}/chats/${peer}/messages`, { body: this.body.trim() })
        this.body = ''
        await this.load()
      }catch(e){ this.errormsg = e.toString() }
    }
  },
  async mounted(){
    await this.load()
    this.timer = setInterval(this.load, 3000)
  },
  beforeUnmount(){ if(this.timer) clearInterval(this.timer) }
}
</script>

<template>
  <div class="container mt-3">
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
    <h4 class="mb-3">Conversation with {{ $route.params.peer }}</h4>
    <div class="border rounded p-3 mb-3" style="height:50vh; overflow:auto; background:#fff;">
      <div v-for="m in msgs" :key="m.id" class="mb-2">
        <small class="text-muted">{{ m.sender }} → {{ m.receiver }} • {{ new Date(m.date).toLocaleString() }}</small>
        <div>{{ m.body }}</div>
      </div>
      <div v-if="msgs.length===0" class="text-muted">No messages yet.</div>
    </div>
    <div class="input-group">
      <input v-model="body" type="text" class="form-control" placeholder="Type a message..." @keyup.enter="send">
      <button class="btn btn-primary" @click="send">Send</button>
    </div>
  </div>
</template>

<style scoped>
</style>


