<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {
	data(){
		return{
			logged: false,
			searchValue: "",
		}
	},
	methods:{
		logout(newValue){
			this.logged = newValue
			this.$router.replace("/login")
		},
		updateLogged(newLogged){
			this.logged = newLogged
		},
		updateView(newRoute){
			this.$router.replace(newRoute)
		},
		search(queryParam){
			this.searchValue= queryParam
			this.$router.replace("/search")
		},
	},

	
	created(){
		if (!localStorage.getItem('notFirstStart')){
			localStorage.clear()
			localStorage.setItem('notFirstStart',true)
			// console.log("first start")
		}
		
	},
	

	mounted(){

		// console.log("Devo modificare ancora lo stile!")
		if (!localStorage.getItem('token')){
			this.$router.replace("/login")
		}else{
			this.logged = true
		}
	},
}
</script>

<template>
    <div class="container-fluid">
		<div class="row">
			<div class="col p-0">
          <main>
					<Navbar v-if="logged" 
					@logoutNavbar="logout" 
					@requestUpdateView="updateView"
					@searchNavbar="search"/>

					<RouterView 
					@updatedLoggedChild="updateLogged" 
					@requestUpdateView="updateView"
					:searchValue="searchValue"/>

            <footer class="py-3 px-4 text-center" style="opacity:.7;">
              <small>WASAgram • a social media sharing UI</small>
            </footer>
				</main>
			</div>
		</div>
	</div>
</template>

<style>
</style>
