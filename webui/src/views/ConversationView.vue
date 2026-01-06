<script>
export default {
  data(){
    return {
      errormsg:null,
      msgs:[],
      body:'',
      loading:false,
      timer:null,
      groupName:'',
      newMember:'',
      groupPhotoPreviewUrl:null
    }
  },
  computed:{
    currentUser(){ return localStorage.getItem('token') },
    peer(){ return `${this.$route.params.peer}` },
    isGroup(){ return this.peer.startsWith('g-') },
    groupId(){
      if(!this.isGroup) return null
      const n = parseInt(this.peer.slice(2), 10)
      return Number.isFinite(n) ? n : null
    }
  },
  methods:{
    async load(){
      try{
        const id = localStorage.getItem('token')
        const peer = encodeURIComponent(this.$route.params.peer)
        const res = await this.$axios.get(`/users/${id}/chats/${peer}/messages`)
        const data = res.data
        this.msgs = Array.isArray(data) ? data : (data && data.messages) ? data.messages : []
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
    },
    async deleteMsg(mid){
      try{
        const id = localStorage.getItem('token')
        const peer = encodeURIComponent(this.$route.params.peer)
        await this.$axios.delete(`/users/${id}/chats/${peer}/messages/${mid}`)
        await this.load()
      }catch(e){ this.errormsg = e.toString() }
    },
    async react(mid){
      try{
        const reaction = window.prompt('Reaction (emoticon/text):')
        if(!reaction) return
        const id = localStorage.getItem('token')
        const peer = encodeURIComponent(this.$route.params.peer)
        await this.$axios.post(`/users/${id}/chats/${peer}/messages/${mid}/comments`, { reaction })
        await this.load()
      }catch(e){ this.errormsg = e.toString() }
    },
    async unreact(mid){
      try{
        const id = localStorage.getItem('token')
        const peer = encodeURIComponent(this.$route.params.peer)
        await this.$axios.delete(`/users/${id}/chats/${peer}/messages/${mid}/comments`)
        await this.load()
      }catch(e){ this.errormsg = e.toString() }
    },
    async forward(mid){
      try{
        const to = window.prompt('Forward to (user id or g-<groupId>):')
        if(!to) return
        const id = localStorage.getItem('token')
        const peer = encodeURIComponent(this.$route.params.peer)
        await this.$axios.post(`/users/${id}/chats/${peer}/messages/${mid}/forward`, { to })
      }catch(e){ this.errormsg = e.toString() }
    },
    async setGroupName(){
      if(!this.isGroup || !this.groupId || !this.groupName.trim()) return
      try{
        await this.$axios.put(`/groups/${this.groupId}`, { name: this.groupName.trim() })
        this.groupName = ''
      }catch(e){ this.errormsg = e.toString() }
    },
    async addMemberToGroup(){
      if(!this.isGroup || !this.groupId || !this.newMember.trim()) return
      try{
        await this.$axios.put(`/groups/${this.groupId}/members/${encodeURIComponent(this.newMember.trim())}`)
        this.newMember = ''
      }catch(e){ this.errormsg = e.toString() }
    },
    async leaveGroup(){
      if(!this.isGroup || !this.groupId) return
      try{
        const me = localStorage.getItem('token')
        await this.$axios.delete(`/groups/${this.groupId}/members/${encodeURIComponent(me)}`)
        this.$router.replace('/chats')
      }catch(e){ this.errormsg = e.toString() }
    },
    async uploadGroupPhoto(){
      if(!this.isGroup || !this.groupId) return
      try{
        let fileInput = document.getElementById('groupPhotoUploader')
        const file = fileInput.files[0]
        if (!file) return
        const reader = new FileReader()
        reader.readAsArrayBuffer(file)
        reader.onload = async () => {
          await this.$axios.put(`/groups/${this.groupId}/photo`, reader.result, { headers: { 'Content-Type': file.type } })
          this.groupPhotoPreviewUrl = URL.createObjectURL(file)
        }
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

    <div v-if="isGroup" class="card mb-3">
      <div class="card-body">
        <h5 class="card-title">Group actions</h5>
        <div class="row g-2 align-items-center mb-2">
          <div class="col-md-5">
            <input v-model="groupName" type="text" class="form-control" placeholder="New group name">
          </div>
          <div class="col-md-2 d-grid">
            <button class="btn btn-outline-secondary" @click="setGroupName" :disabled="!groupName.trim()">Set name</button>
          </div>
          <div class="col-md-5">
            <div class="input-group">
              <input v-model="newMember" type="text" class="form-control" placeholder="Member user id">
              <button class="btn btn-outline-secondary" @click="addMemberToGroup" :disabled="!newMember.trim()">Add</button>
            </div>
          </div>
        </div>

        <div class="row g-2 align-items-center">
          <div class="col-md-6">
            <input id="groupPhotoUploader" type="file" accept="image/png,image/jpeg" class="form-control">
          </div>
          <div class="col-md-3 d-grid">
            <button class="btn btn-outline-secondary" @click="uploadGroupPhoto">Upload photo</button>
          </div>
          <div class="col-md-3 d-grid">
            <button class="btn btn-outline-danger" @click="leaveGroup">Leave group</button>
          </div>
        </div>
        <img v-if="groupPhotoPreviewUrl" :src="groupPhotoPreviewUrl" alt="Group photo preview" style="max-width: 140px; max-height: 140px;" class="mt-2">
      </div>
    </div>

    <div class="border rounded p-3 mb-3" style="height:50vh; overflow:auto; background:#fff;">
      <div v-for="m in msgs" :key="m.id" class="mb-2">
        <small class="text-muted">
          {{ m.sender }} → {{ m.receiver }} • {{ new Date(m.date).toLocaleString() }}
          <span v-if="m.sender === localStorage.getItem('token')">
            <span v-if="m.status === 2"> ✓✓</span>
            <span v-else-if="m.status === 1"> ✓</span>
          </span>
        </small>
        <div>{{ m.body }}</div>
        <div v-if="m.reactions && m.reactions.length" class="small text-muted">
          Reactions:
          <span v-for="(r,i) in m.reactions" :key="i">{{ r.reaction }} ({{ r.user_id || r.userId }}) </span>
        </div>
        <div class="mt-1">
          <button class="btn btn-sm btn-outline-secondary me-1" @click="react(m.id)">React</button>
          <button class="btn btn-sm btn-outline-secondary me-1" @click="unreact(m.id)">Unreact</button>
          <button class="btn btn-sm btn-outline-secondary me-1" @click="forward(m.id)">Forward</button>
          <button v-if="m.sender === currentUser" class="btn btn-sm btn-outline-danger" @click="deleteMsg(m.id)">Delete</button>
        </div>
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


