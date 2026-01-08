<script>
export default {
  data(){
    return { errormsg: null, peers: [], groupName: '', groupMembers: '' }
  },
  methods:{
    async load(){
      try{
        const id = localStorage.getItem('token')
        const res = await this.$axios.get(`/users/${id}/chats`)
        const data = res.data
        this.peers = Array.isArray(data) ? data : (data && data.conversations) ? data.conversations : []
      }catch(e){ this.errormsg = e.toString() }
    },
    async createGroup(){
      try{
        this.errormsg = null
        const id = localStorage.getItem('token')
        const members = (this.groupMembers || '')
          .split(',')
          .map(s => s.trim())
          .filter(s => s.length > 0)
        await this.$axios.post(`/users/${id}/groups`, { name: this.groupName.trim(), members })
        this.groupName = ''
        this.groupMembers = ''
        await this.load()
      }catch(e){ this.errormsg = e.toString() }
    },
    open(peer){
      const id = peer && (peer.peer || peer.user_id || peer.id_user || peer.IdUser || peer.id || `${peer}`)
      this.$router.push(`/chats/${encodeURIComponent(id)}`)
    }
  },
  async mounted(){ await this.load() }
}
</script>

<template>
  <div class="container mt-4">
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
    <h3 class="mb-3">Chats</h3>
    <div class="card mb-3">
      <div class="card-body">
        <h5 class="card-title mb-3">Create group</h5>
        <div class="row g-2">
          <div class="col-md-4">
            <input v-model="groupName" type="text" class="form-control" placeholder="Group name">
          </div>
          <div class="col-md-6">
            <input v-model="groupMembers" type="text" class="form-control" placeholder="Members (comma-separated user ids)">
          </div>
          <div class="col-md-2 d-grid">
            <button class="btn btn-outline-primary" @click="createGroup" :disabled="!groupName.trim()">Create</button>
          </div>
        </div>
      </div>
    </div>
    <div v-if="peers.length===0" class="text-muted">No conversations yet.</div>
    <ul class="list-group">
      <li v-for="(u,i) in peers" :key="i" class="list-group-item d-flex justify-content-between align-items-center" @click="open(u)">
        <div class="d-flex flex-column">
          <strong>{{ u.name || u.nickname || u.peer || u.IdUser || u.id_user }}</strong>
          <small class="text-muted">
            {{ u.lastMessagePreview || '' }}
            <span v-if="u.lastMessageAt"> • {{ new Date(u.lastMessageAt).toLocaleString() }}</span>
          </small>
        </div>
        <button class="btn btn-sm btn-primary">Open</button>
      </li>
    </ul>
  </div>
  
</template>

<style scoped>
.list-group-item{ cursor:pointer }
</style>


