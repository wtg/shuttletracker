<template>
	<div id='select'></div>
</template>

<script>
const Tabulator = require('tabulator-tables');
export default {
  data() {
    return {
      tabulator: null,
      tableColumn: [
        {title: 'Day', field: 'day', align: 'center', headerSort:false, cellClick: function(e, cell){console.log(cell.getValue()); console.log(cell.getColumn().getField())},}, 
        {title: 'Time', field: 'time', align: 'center', headerSort:false, cellClick: function(e, cell){console.log(cell.getValue()); console.log(cell.getColumn().getField())},}, 
      ],
    };
  },

  mounted() {
    this.tabulator = new Tabulator('#select', {
      data: [],
      columns: this.tableColumn,
      height: 250,
      layout: 'fitColumns',
      placeholder: 'Select Times',
    });
  },

  watch: {
    tableData: {
      handler: function(newData) {
        this.tabulator.replaceData(newData);
      },
      deep: true,
    },
  },
};
</script>

<style>
#select {
  background-color:#666;
  border: 1px solid #333;
  border-radius: 10px;
}

#select .tabulator-header {
  background-color:#666;
  color:#fff;
}

#select .tabulator-header .tabulator-col,
#select .tabulator-header .tabulator-col-row-handle {
  white-space: normal;
  background-color:#333;
}

#select .tabulator-tableHolder .tabulator-table .tabulator-row {
  background-color:#666;
  color:#fff;
}

#select .tabulator-tableHolder .tabulator-table .tabulator-row:nth-child(even) {
  background-color:#444;
}
</style>