<script setup lang="ts">
import 'leaflet/dist/leaflet.css';
import 'leaflet.markercluster/dist/MarkerCluster.css';
import 'leaflet.markercluster/dist/MarkerCluster.Default.css';
import L from 'leaflet';
import 'leaflet.markercluster';
import {onMounted, ref} from 'vue'
import {ElNotification as nofity} from 'element-plus'
import apiUrl from "@/config/config.ts";
import apiClient from "@/axios";

const checkLine = ref(false)
let map = ref<L.Map | null>(null);
let ponitList = [Array];
let polyline: L.Polyline | null = null;
const markers = ref<L.CircleMarker[]>([]); // 存储绘制的点


// 格式化日期为字符串（YYYY-MM-DD HH:mm:ss）
const formatDateTime = (date, isEndTime = false) => {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0'); // 月份从 0 开始
  const day = String(date.getDate()).padStart(2, '0');

  // 如果是 endTime，手动设置时间为 00:00:00
  const hours = isEndTime ? '00' : String(date.getHours()).padStart(2, '0');
  const minutes = isEndTime ? '00' : String(date.getMinutes()).padStart(2, '0');
  const seconds = isEndTime ? '00' : String(date.getSeconds()).padStart(2, '0');

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
};

const initMap = async () => {
  map.value = L.map('map', {
    center: [29.44916482692468, 106.47399902343751],//中心坐标
    zoom: 4,//初始缩放，因为在下文写了展示全地图，所以这里不设置，也可以设置
    minZoom: 1,
    maxZoom: 20,
    zoomControl: true, //缩放组件
    // attributionControl: false, //去掉右下角logol
  })
// let map=L.map("map").setview([29.44916482692468, 106.47399902343751],9);


  L.tileLayer(
      "https://webrd04.is.autonavi.com/appmaptile?lang=zh_cn&size=1&scale=1&style=7&x={x}&y={y}&z={z}",
      {
        attribution:
            '© <p>OpenStreetMap</p> contributors',
      }
  ).addTo(map.value!);

  // L.tileLayer('http://webrd0{s}.is.autonavi.com/appmaptile?lang=zh_cn&size=1&scale=1&style=8&x={x}&y={y}&z={z}', {
  //   subdomains: ['1', '2', '3', '4'],
  //   minZoom: 1,
  //   maxZoom: 19
  // }).addTo(map);

//   let markerIcon = L.icon({
//     iconUrl: "https://unpkg.com/leaflet@1.9.3/dist/images/marker-icon.png",
//     iconSize: [20, 30],
//   });
//   let marker = L.marker([29.44916482692468, 106.47399902343751], { style: 'red' })
//       .addTo(map)
//       .bindPopup("标记")
//       .openPopup();
// }

// 获取当前日期
  const now = new Date();

// 计算 beginTime（当前日期 - 3 天）
  const beginTime = new Date(now);
  beginTime.setDate(now.getDate() - 3);

// 获取格式化后的 beginTime 和 endTime
  const formattedBeginTime = formatDateTime(beginTime);
  const formattedEndTime = formatDateTime(now, true); // 设置 isEndTime 为 true

  queryDate.value = []
  queryDate.value.push(formattedBeginTime)
  queryDate.value.push(formattedEndTime)

// 配置请求头
  const config = {
    headers: {
      'sessionID': sessionStorage.getItem('session'), // 假设你的session信息是以这种方式传递的
      // 你可以在这里添加更多的请求头
    },
    params: {
      beginTime: formattedBeginTime, // 添加 beginTime 查询参数
      endTime: formattedEndTime // 添加 endTime 查询参数
    }
  };


  const response = await apiClient.get('/getData', config);
  if (!response?.data.success) {
    nofity.error({
      message: response?.data.msg,
    })
    return
  }

  if (response.data.data === null) {
    nofity.error({
      message: '所选日期无数据',
    })
    return
  }
  ponitList = response.data.data

  //无聚合
  response.data.data.forEach(function(coord) {
    const latLngString = `经纬度: ${coord[0].toFixed(4)}, ${coord[1].toFixed(4)}`;
    const marker = L.circleMarker(coord, {
      color: 'red',
      radius: 1 // 设置半径
    }).addTo(map.value!);
    markers.value.push(marker); // 保存新绘制的点
  });



  // // 创建聚合器
  // const markers = L.markerClusterGroup({
  //   singleMarkerMode: true,//true:单个marker显示聚合数字1,false:显示单个marker
  //   maxClusterRadius: 5, // 设置最大聚合半径
  //   showCoverageOnHover: false, // 鼠标悬停时显示覆盖区域
  //   zoomToBoundsOnClick: true // 点击聚合标记时缩放到包含的所有标记
  // });
  //
  // response.data.data.forEach(function(coord) {
  //   const latLngString = `经纬度: ${coord[0].toFixed(4)}, ${coord[1].toFixed(4)}`;
  //   // L.circleMarker(coord, {
  //   //   color: 'red',
  //   //   radius: 1 // 设置半径
  //   // }).addTo(map).bindPopup(latLngString).openPopup();
  //
  //   const marker = L.circleMarker(coord, {
  //     color: 'red',
  //     radius: 2
  //   }).bindPopup(latLngString).openPopup();
  //
  //   // const marker = L.marker(L.latLng(coord[0], coord[1])).bindPopup(latLngString).openPopup();
  //
  //   markers.addLayer(marker);
  // });
  // // 将聚合器添加到地图
  // map.addLayer(markers);


}

const setLine = (value: boolean) => {
  console.log(value)
  if (value && map.value) {
    //创建轨迹连线
    polyline = L.polyline(ponitList, {
      color: 'red', // 线的颜色
      weight: 3, // 线的宽度
      opacity: 0.5, // 线的透明度
    }).addTo(map.value!)
    map.value!.fitBounds(polyline.getBounds());
  } else {
    map.value.removeLayer(polyline)
  }
}

const size = ref<'default' | 'large' | 'small'>('default')
const queryDate = ref([])

const shortcuts = [
  {
    text: 'Last week',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
      return [start, end]
    },
  },
  {
    text: 'Last month',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
      return [start, end]
    },
  },
  {
    text: 'Last 3 months',
    value: () => {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
      return [start, end]
    },
  },
]

const queryData = async () => {
  console.log(queryDate.value)
  // 配置请求头
  const config = {
    headers: {
      'sessionID': sessionStorage.getItem('session'), // 假设你的session信息是以这种方式传递的
      // 你可以在这里添加更多的请求头
    },
    params: {
      beginTime: queryDate.value[0], // 添加 beginTime 查询参数
      endTime: queryDate.value[1] // 添加 endTime 查询参数
    }
  };
  const response = await apiClient.get('/getData', config);
  if (!response?.data.success) {
    nofity.error({
      message: response?.data.msg,
    })
    return
  }
  if (response.data.data === null) {
    nofity.error({
      message: '所选日期无数据',
    })
    return
  }
  ponitList = response.data.data

  // 清除之前的点
  markers.value.forEach(marker => {
    map.value!.removeLayer(marker); // 从地图上移除
  });
  markers.value = []; // 清空 markers 数组

  //无聚合
  response.data.data.forEach(function (coord) {
    const latLngString = `经纬度: ${coord[0].toFixed(4)}, ${coord[1].toFixed(4)}`;
    const marker = L.circleMarker(coord, {
      color: 'red',
      radius: 1 // 设置半径
    }).addTo(map.value!);
    markers.value.push(marker); // 保存新绘制的点
  });
}

const exportCSV = async () => {
  try {
    const config = {
      headers: {
        'sessionID': sessionStorage.getItem('session'), // 假设你的session信息是以这种方式传递的
        // 你可以在这里添加更多的请求头
      },
      params: {
        beginTime: queryDate.value[0], // 添加 beginTime 查询参数
        endTime: queryDate.value[1] // 添加 endTime 查询参数
      },
      // responseType: 'blob'
    };
    const response = await apiClient.get('/export', config);

    // 检查响应的 Content-Type
    const contentType = response.headers['content-type'];
    const contentDisposition = response.headers['content-disposition'];

    // 判断是返回 JSON 数据还是 CSV 文件
    if (contentType === 'application/json') {
      const result = await response.data;
      console.log(result)
      if (result?.msg) {
        nofity.error({
          message: result?.msg,
        })
        return
      }
      if (result.flag === 'client') {
        // 处理小于 1 万条的数据
        const data = result.data;
        if (data === null) {
          nofity.error({
            message: '所选日期无数据',
          })
          return
        }
        exportToCSV(data);
      }
    } else if (contentType === 'text/csv') {
      // 处理大于 1 万条的数据，自动下载 CSV 文件
      const blob = new Blob([response.data], { type: 'text/csv' });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      const filename = contentDisposition
          ? contentDisposition.split('filename=')[1].replace(/"/g, '')
          : 'footprint_data.csv';
      a.href = url;
      a.setAttribute('download', filename);
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
    }
  } catch (error) {
    console.error('Error fetching data:', error);
  }
}

function exportToCSV(data: any) {
  // CSV 文件头
  const header = [
    "geoTime", "latitude", "longitude", "altitude", "course",
    "horizontalAccuracy", "verticalAccuracy", "speed",
    "status", "activity", "network", "appStatus",
    "dayTime", "groupTime", "isSplit", "isMerge",
    "isAdd", "networkName"
  ];

  // 将数据转换为 CSV 格式
  const csvContent = [header, ...data].map(row => row.join(",")).join("\n");
  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.setAttribute('download', 'data.csv');
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
}

onMounted(() => {
  initMap()
})
</script>

<template>
  <div style="display: flex;flex-direction: row">
    <div id="map"></div>
    <div>
      <el-checkbox label="连线" style="margin-left: 20px" v-model="checkLine" @change="setLine" v-if="false"/>
      <el-date-picker
          v-model="queryDate"
          type="daterange"
          unlink-panels
          range-separator="To"
          start-placeholder="Start date"
          end-placeholder="End date"
          :shortcuts="shortcuts"
          :size="size"
          value-format="YYYY-MM-DD"
      />
      <el-button @click="queryData">查询</el-button>
      <el-button @click="exportCSV">导出为csv</el-button>
    </div>
  </div>
</template>

<style>
#map {
  height: 800px;
  width: 800px;
  margin: auto;
  border: 10px solid pink;
  overflow: auto;
}
</style>
