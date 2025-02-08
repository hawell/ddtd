<template>
  <v-card
    color="grey-lighten-4"
    rounded="0"
    flat
  >
    <v-toolbar collapse>
      <v-btn icon @click="refresh">
        <v-icon>mdi-refresh</v-icon>
      </v-btn>
      <v-btn icon @click="clear">
        <v-icon>mdi-delete</v-icon>
      </v-btn>
    </v-toolbar>
  </v-card>
  <v-container class="ma-0 pa-0" style="max-width: 100%;">
    <v-row class="ma-0 pa-0">
      <v-col :cols="10">
        <div class="my-2" id="timeline-container" ref="timelineContainer" style="height: 60dvh;">
        </div>
        <json-tree-view class="scrollable elevation-2" style="height: 30dvh;" :json="traceValue">
        </json-tree-view>
      </v-col>
      <v-col :cols="2">
        <v-list class="elevation-2" style="height: 100%;">
          <v-list-item v-for="item in parentTraces" :key="item.id" :value="item" @click="parentSelected(item)">
            <v-list-item-title v-text="item.text"></v-list-item-title>
            <v-list-item-subtitle>{{ item.time }}</v-list-item-subtitle>
          </v-list-item>
        </v-list>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { createVuetify } from 'vuetify'
import { Timeline, DataSet } from "vis-timeline/standalone"
import { onMounted } from "vue"
import { JsonTreeView } from "json-tree-view-vue3";
import 'json-tree-view-vue3/dist/style.css'

var timeline = ref(null)
const timelineContainer = ref(null)
const traceValue = ref("{}")
const parentTraces = ref([])

// Create a DataSet (allows two way data-binding)
var items = new DataSet([]);
var groups = new DataSet([
{
  id: 10,
  title: "traces",
  content: "traces",
  treeLevel: 1,
}
]);

// Configuration for the Timeline
var options = {
  minHeight: "100%",
  maxHeight: "100%",
  margin: { item: 0 },
  order: customOrder,
};

onMounted(() => {
  var container = document.getElementById("timeline-container")
  timeline = new Timeline(container, items, groups, options)
  timeline.on('select', onSelect)
})

function refresh() {
  $fetch('http://localhost:8866/api/get_traces').then(data => {
    //console.debug(data)
    items.clear()
    traceValue.value = "{}"
    parentTraces.value = []
    var end = 0
    var start = 4000000000000
    data.forEach(trace => {
      const item = {}
      item.name = trace.name
      item.content = trace.name
      item.start = trace.start / 1000000
      item.end = item.start + trace.duration / 1000000
      item.group = 10
      item.trace = trace
      item.level = item.trace.level
      item.id = trace.span_id
      items.add(item)
      if (trace.parent_id == 0) {
        const start = new Date(item.start)
        parentTraces.value.push({id: item.id, text: item.name, time: start.toISOString(start)})
      }
      if (item.start < start) {
        start = item.start
      }
      if (item.end > end) {
        end = item.end
      }
    });
    //console.debug(data, start, end)
    timeline.setOptions({"start": start, "end": end})
  })
}

async function clear() {
  await $fetch('http://localhost:8866/api/clear', { method: 'POST'})
  traceValue.value = "{}"
  parentTraces.value = []
  items.clear()
}

function onSelect(properties) {
  //console.debug(properties)
  if (properties.items.length == 0) {
    traceValue.value = "{}"
    return
  }
  const item = items.get(properties?.items[0])
  traceValue.value = JSON.stringify(item?.trace)
}

function customOrder(a, b) {
  // order by id
  return a.level - b.level
}

function parentSelected(trace){
  //console.debug(trace.id)
  traceValue.value = "{}"
  const item = items.get(trace.id)
  timeline.setOptions({"start": item.start, "end": item.end})

}

</script>

<style scoped>
html{
    height: 100%;
}

body{
    margin:0;
    padding:0;
    overflow:hidden;
    height:100%;
}

#container{
    width:1000px;
    margin:0 auto;
    height:100%;
}
.scrollable {
   overflow-y: scroll;
}
</style>