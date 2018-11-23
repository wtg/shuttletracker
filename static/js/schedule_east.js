 Vue.component('demo-grid', {
            template: '#grid-template',
            props: {
                rows: Array,
                columns: Array,
                filterKey: false,			
            },
            data:
				function() {
                var colOrders = {}
                this.columns.forEach(function(col) {
                    colOrders[col] = 1
                })
                return {
                    setOrder: colOrders,

                }

            },
            computed: {
                filterData: function() {
                    var data = this.rows.slice()
					var index = 0
                    var filterKey = this.filterKey
					var offset = 4

				//Filter Key is a trigger
                   if (filterKey) {
					   //Place filter
					   if (list.place_condition.length > 0){
						data = data.filter(function(row) {	
						for(var k=0; k < list.place_condition.length; k++){
							if (row.name === list.place_condition[k]) {
								return true;
							}
							if (k == list.place_condition.length - 1){
								return false;
							}
						}							    
                      })
					}
					//Time filter
					if (list.time_condition.length > 0){
						data = data.filter(function(row) {	
						for(var k=0; k < list.time_condition.length; k++){
							if (String(row.bus1).indexOf(list.time_condition[k]) > -1) {
								return true;
							}
							if (k == list.time_condition.length - 1){
								return false;
							}
						}							    
                      })
					}
								
                    }

                    return data
                }
            },
            methods: {
				//Make the table order
                order: function(key) {
                    this.setOrder[key] = this.setOrder[key] * -1
                },
            }

        })
        var vm = new Vue({
            el: '#grid',
            data: {
                columnData: ['name', 'bus1','bus2'],
                searchQuery: '',
				list,
                rowData: [{
                    name: 'Union',
                    bus1: '7:00a',
					bus2: '7:20a'
                }, {
                    name: 'Tibbits Ave',
                    bus1: '7:03a',
					bus2: '7:23a'
                },{
                    name: 'B-Lot',
                    bus1: '7:06a',
					bus2: '7:26a'
                }, {
                    name: 'Colonie',
                    bus1: '7:09a',
					bus2: '7:29a'
                }, {
                    name: 'Union',
                    bus1: '10:19a',
					bus2: '10:04a'
                },{
                    name: 'Tibbits Ave',
                    bus1: '10:03a',
					bus2: '10:23a'
                },{
                    name: 'B-Lot',
                    bus1: '10:06a',
					bus2: '10:26a'
                }, {
                    name: 'Colonie',
                    bus1: '10:09a',
					bus2: '10:29a'
                }]
            },

        })
		
		
		var list={
	category:[
		{
			name:'Stop',
			items:[
			{
				name:'Union',
				value:'Union',
				active: false
			},
			{
				name:'Tibbits Ave',
				value:'Tibbits Ave',
				active: false
			},
			{
				name:'B-Lot',
				value:'B-Lot',
				active: false
			},
			{
				name:'Colonie',
				value:'Colonie',
				active: false
			}
			]
		},
		{
			name:'Time',
			items:[
			{
				name:'7:00',
				value:'7:',
				active: false
			},
			{
				name:'8:00',
				value:'8:',
				active: false
			},
			{
				name:'9:00',
				value:'9:',
				active: false
			},
			{
				name:'10:00',
				value:'10:',
				active: false
			},
			{
				name:'11:00',
				value:'11:',
				active: false
			
			},
			{
				name:'12:00',
				value:'12:',
				active: false
			
			},
			{
				name:'13:00',
				value:'13:',
				active: false
			
			},
			{
				name:'14:00',
				value:'14:',
				active: false
			
			},
			{
				name:'15:00',
				value:'15:',
				active: false
			
			},
			{
				name:'16:00',
				value:'16:',
				active: false
			
			},
			{
				name:'17:00',
				value:'17:',
				active: false
			
			},
			{
				name:'18:00',
				value:'18:',
				active: false
			
			},
			{
				name:'19:00',
				value:'19:',
				active: false
			
			},
			{
				name:'20:00',
				value:'20:',
				active: false
			
			},
			{
				name:'21:00',
				value:'21:',
				active: false
			
			},
			{
				name:'22:00',
				value:'22:',
				active: false
			
			}
			]
		},
		
		
	],
	condition:[ //Conditions shows to users
	],
	place_condition:[ //Place conditions that users choose

	],
	time_condition:[//Time conditions that users choose

	]
};
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
							this.list.place_condition.push(v.value)
						}						
					}
					else{//When index is 1, uses choose time. Add condition to place condiiton
						if (!this.list.time_condition.includes(v.value)){
							this.list.time_condition.push(v.value)
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
					this.list.place_condition = []
					
					for(var k=0; k < item.length; k++){
						this.list.place_condition.push(item[k].value)
					}
				}
				else {//When index is 0, uses choose place. Add all conditions in place condiiton
					this.list.time_condition = []
					for(var k=0; k < item.length; k++){
						this.list.time_condition.push(item[k].value)
					}
				}
				
			});	
		
				
		}
	}
});