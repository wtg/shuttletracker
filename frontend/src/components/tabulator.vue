<template>
	<div id="table"></div>
</template>

<script>
//npm install tabulator-tables --save
const Tabulator = require("tabulator-tables");
export default {
	data() {
		return {
			tabulator: null,
			tableColumn: [
				{title:"Sunday",field:"su",width:100,align:"center"},
				{title:"Monday",field:"m",width:100,align:"center"},
				{title:"Tuesday",field:"t",width:100,align:"center"},
				{title:"Wednesday",field:"w",width:100,align:"center"},
				{title:"Thurday",field:"th",width:100,align:"center"},
				{title:"Friday",field:"f",width:100,align:"center"},
				{title:"Saturday",field:"sa",width:100,align:"center"},
			],
			tableData: [
				{id:1,su:"8:00",m:"8:00",t:"8:00",w:"8:00",th:"8:00",f:"8:00",sa:"8:00"},
				{id:2,su:"9:00",m:"9:00",t:"9:00",w:"9:00",th:"9:00",f:"9:00",sa:"9:00"},
				{id:3,su:"10:00",m:"10:00",t:"10:00",w:"10:00",th:"10:00",f:"10:00",sa:"10:00"},
				{id:4,su:"11:00",m:"11:00",t:"11:00",w:"11:00",th:"11:00",f:"11:00",sa:"11:00"},
				{id:5,su:"12:00",m:"12:00",t:"12:00",w:"12:00",th:"12:00",f:"12:00",sa:"12:00"},
				{id:6,su:"13:00",m:"13:00",t:"13:00",w:"13:00",th:"13:00",f:"13:00",sa:"13:00"},
				{id:7,su:"14:00",m:"14:00",t:"14:00",w:"14:00",th:"14:00",f:"14:00",sa:"14:00"},
				{id:8,su:"15:00",m:"15:00",t:"15:00",w:"15:00",th:"15:00",f:"15:00",sa:"15:00"},
				{id:9,su:"16:00",m:"16:00",t:"16:00",w:"16:00",th:"16:00",f:"16:00",sa:"16:00"},
				{id:10,su:"17:00",m:"17:00",t:"17:00",w:"17:00",th:"17:00",f:"17:00",sa:"17:00"},
				{id:11,su:"18:00",m:"18:00",t:"18:00",w:"18:00",th:"18:00",f:"18:00",sa:"18:00"},
			],
			addData: {id:null,su:null,m:null,t:null,w:null,th:null,f:null,sa:null}
		}
	},
	mounted(){
		this.tabulator = new Tabulator("#table", {
			data: this.tableData,
			columns: this.tableColumn,
			height: 500,
			layout:"fitColumns",
		});
	},
	methods: {
		addTabulator(){
			const obj = Object.assign({},this.addData)
			obj.id = this.tableData.length
			this.tableData.push(obj)
			Object.keys(this.addData).forEach(key => {
				this.addData[key] = null
			})
		}
	},
	watch: {
		tableData:{
			handler: function (newData) {
				this.tabulator.replaceData(newData);
			},
			deep: true,
		}
	}
}
</script>

<style>
@import '~tabulator-tables/dist/css/tabulator.min.css'
</style>