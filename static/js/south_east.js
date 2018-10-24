 Vue.component('demo-grid', {
            template: '#grid-template',
            props: {
                rows: Array,
                columns: Array,
                filterKey: String
            },
            data: function() {
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
                    var filterKey = this.filterKey && this.filterKey.toLowerCase()
                    if (filterKey) {
                        data = data.filter(function(row) {
							
                            return Object.keys(row).some(function(key) {
							
                                return String(row[key]).toLowerCase().indexOf(filterKey) > -1
                            })
                        })
                    }
                    return data
                }
            },
            methods: {
                //对一个包含对象的数组的排序，需要提供一个对象键并以此值来进行排序
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
                }]
            },

        })