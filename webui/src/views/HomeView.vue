<script>
export default {
	data: function () {
		return {
			errormsg: null,
			photos: [],
		}
	},

	methods: {
		
		async loadStream() {
			try {
				this.errormsg = null
				// Home get: "/users/:id/home"
				let response = await this.$axios.get("/users/" + localStorage.getItem('token') + "/home")

				if (response.data != null){
					this.photos = response.data
				}
				// Fallback to own posts if stream is empty
				if (!this.photos || this.photos.length === 0){
					try{
						const prof = await this.$axios.get("/users/"+localStorage.getItem('token'))
						this.photos = prof.data && prof.data.posts ? prof.data.posts : []
					}catch(_){/* ignore */}
				}
			} catch (e) {
				this.errormsg = e.toString()
			}
		}
	},

	async mounted() {
		await this.loadStream()
    // Listen for global refresh requests (e.g., after upload)
    this._onStreamRefresh = () => { this.loadStream() }
    window.addEventListener('stream:refresh', this._onStreamRefresh)
    // Prepend newly uploaded photo without refetch
    this._onNewPhoto = (e) => {
      if (e && e.detail){
        this.photos.unshift(e.detail)
      }
    }
    window.addEventListener('stream:new-photo', this._onNewPhoto)
  },
  computed: {
    currentUserId(){
      return localStorage.getItem('token')
    }
  },

  beforeUnmount(){
    if (this._onStreamRefresh){
      window.removeEventListener('stream:refresh', this._onStreamRefresh)
    }
    if (this._onNewPhoto){
      window.removeEventListener('stream:new-photo', this._onNewPhoto)
    }
	}

}
</script>

<template>
	<div class="container-fluid">
		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>

		<div class="row">
			<Photo
				v-for="(photo,index) in photos"
				:key="index"
				:owner="photo.owner"
				:photo_id="photo.photo_id"
                :comments="photo.comments != null ? photo.comments : []"
                :likes="photo.likes != null ? photo.likes : []"
				:upload_date="photo.date"
                :isOwner="photo.owner === currentUserId"
			/>
		</div>

		<div v-if="photos.length === 0" class="row ">
			<h1 class="d-flex justify-content-center mt-5" style="color: white;">There's no content yet, follow somebody!</h1>
		</div>
	</div>
</template>

<style>
</style>
