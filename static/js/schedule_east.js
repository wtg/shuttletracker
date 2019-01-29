 Vue.component('demo-grid', {
    template: '#grid-template',
    props: {
        rows: Array,
        columns: Array,
        filterKey: false,		
    },
    data:
		function() {
            var colOrders = {};
            this.columns.forEach(function(col) {
                colOrders[col] = 1;
            })
            return {
                setOrder: colOrders,
            }
        },
    computed: {
        filterData: function() {
            var data = this.rows.slice(); //hold the data of lists
            var filterKey = this.filterKey;
			//Filter Key is the trigger of searching
            if (filterKey) {
				//Filter
				data = data.filter(function(row) {	
				//Place filter
				for(var k=0; k < list.place_condition.length; k++){
					if (row.name === list.place_condition[k]) {
						if (list.time_condition.length > 0){
							//Time filter
							for(var j=0; j < list.time_condition.length; j++){
								if (row.stop1.match(list.time_condition[j])) {
								
									return true;	
								}			
							}
							return false;
						}	
					}
				}	
				return false;						
				})		
            }
            return data;
        }
    },
    methods: {
		//Make the table order
        order: function(key) {
             this.setOrder[key] = this.setOrder[key] * -1;
            },
        }

    })

var list={
	category:[],
	condition:[ //Conditions shows to users
	],
	place_condition:[ //Place conditions that users choose

	],
	time_condition:[//Time conditions that users choose

	]
};
var vm = new Vue({
    el: '#grid',
    data: {
        columnData: ['name', 'stop1','stop2'], //Show the length and content of column
        searchQuery: '',
		list,
		rowData : [],
    },
	created (){ //Fetch JSON file
		fetch("/static/json/east_category.json")
		.then (response => response.json())
		.then (json => {
		this.list.category = json.category
		});
		fetch("/static/json/east_normal.json")
		.then (response => response.json())
		.then (json => {
		this.rowData = json.rowData
		});
		
	}
})
			

var count=0;
var app =new Vue({
	el:'#app',
	data:list,
	methods:{
		//Add conditions
		handle:function(index,key){
			var item=this.category[index].items;
			item.filter(function(v,i){
				//Add contents into condition
				if(i==key){
					v.active=true;			
					this.list.condition.filter(function(v2,i){
						if(v.name==v2.name){
							this.list.condition.splice(i,1);
							count--;
						}
					});					
					Vue.set(this.list.condition,count++,v);
					if (index == 0){ //When index is 0, uses choose place. Add condition to place condiiton
						if (!this.list.place_condition.includes(v.value)){
							this.list.place_condition.push(v.value);
						}						
					}
					else{//When index is 1, uses choose time. Add condition to place condiiton
						if (!this.list.time_condition.includes(v.value)){
							this.list.time_condition.push(v.value);
						}	
					}
						
				}
			});
			
		},
		//Remove conditions
		remove:function(index){
			var item=this.category[index].items;
		
			//Clear conditions
			item.filter(function (v1,key) {
				v1.active=false;
				this.list.condition.filter(function(v2,i){
					if(v1.name==v2.name){
						this.list.condition.splice(i,1);
						count--;
					}
				});
				if (index == 0){//When index is 0, uses choose place. Clear conditions in place condiiton
					this.list.place_condition = [];
				}
				else {//When index is 0, uses choose place. Clear conditions in place condiiton
					this.list.time_condition = [];
				}
			});
		},
		//Add all the conditions into condition
		allIn:function(index){
			var item=this.category[index].items;
			item.filter(function (v,key) {
				v.active=true;
				this.list.condition.filter(function(v2,i){
					if(v.name==v2.name){
						this.list.condition.splice(i,1);
						count--;
					}
				});					
				Vue.set(this.list.condition,count++,v);
				if (index == 0){//When index is 0, uses choose place. Add all conditions in place condiiton
					this.list.place_condition = [];	
					for(var k=0; k < item.length; k++){
						this.list.place_condition.push(item[k].value);
					}
				}
				else {//When index is 0, uses choose place. Add all conditions in place condiiton
					this.list.time_condition = [];
					for(var k=0; k < item.length; k++){
						this.list.time_condition.push(item[k].value);
					}
				}
				
			});	
		
				
		},
		
	}
});