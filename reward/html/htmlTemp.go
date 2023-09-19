package html

const HtmlTemplate = `<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>WebUI</title>
    <!-- import Element-UI CSS -->
    <link rel="stylesheet" href="https://unpkg.com/element-ui/lib/theme-chalk/index.css">
</head>

<body>
    <div id="app">
        <el-container>
            <el-header>
                <h1>WebUI</h1>
            </el-header>
            <el-main>
                <el-form ref="form" :model="settings" label-width="100px">
                    <el-form-item label="Proxy On">
                        <el-switch v-model="settings.proxy_on"></el-switch>
                    </el-form-item>
                    <el-form-item label="Proxy Address">
                        <el-input v-model="settings.proxy" @blur="saveSettings"></el-input>
                    </el-form-item>
                    <el-form-item label="Key Words">
                        <el-tag v-for="(word, index) in settings.key_words" :key="index" closable
                            @close="removeKeyword(index)">
                            {{ word }}
                        </el-tag>
                        <el-input v-model="newKeyword" placeholder="Enter a new keyword"
                            @keyup.enter="addKeyword"></el-input>
                        <el-button @click="addKeyword" type="success">Add</el-button>
                    </el-form-item>

                    <el-form-item label="Cookie">
                        <el-input type="textarea" :rows="5" placeholder="键入cookies"
                            v-model="settings.cookies" @blur="saveSettings"></el-input>
                    </el-form-item>

                    <el-button @click="saveSettings" type="primary">Save</el-button>
                    <el-button @click="startGetPoints" type="primary">开始</el-button>
                </el-form>
            </el-main>
            <el-header>
                <h3>output</h3>
            </el-header>
<!--            <ul>-->
<!--                {{range .Data}}-->
<!--                <li>{{.}}</li>-->
<!--                {{end}}-->
<!--            </ul>-->
            刷分进度: <span id="process">{{.}}</span>
        </el-container>
    </div>

    <!-- import Vue and Element-UI -->
    <script src="https://unpkg.com/vue@2/dist/vue.js"></script>
    <script src="https://unpkg.com/element-ui/lib/index.js"></script>
    <!-- import axios -->
    <script src="https://unpkg.com/axios/dist/axios.min.js"></script>
<!--    <script>-->
<!--        var socket = new WebSocket("ws://{{.Host}}/getInfo");-->

<!--        socket.onmessage = function(event) {-->
<!--            var newData = JSON.parse(event.data);-->
<!--            var list = document.querySelector("ul");-->
<!--            list.innerHTML = "";-->

<!--            for (var i = 0; i < newData.length; i++) {-->
<!--                var li = document.createElement("li");-->
<!--                li.textContent = newData[i];-->
<!--                list.appendChild(li);-->
<!--            }-->
<!--        };-->
<!--    </script>-->
    <script>
        var socket = new WebSocket("ws://{{.}}/getInfo");

        socket.onmessage = function(event) {
            document.getElementById("process").textContent = event.data;
        };
    </script>
    <script>
        new Vue({
            el: '#app',
            data: {
                settings: {
                    proxy_on: false,
                    proxy: '',
                    key_words: [],
                    cookies: '' // New cookie textarea field
                },
                keywords: ['关键字1', '关键字2', '关键字3'],
                newKeyword: '', // New keyword to be added
                headers : {
                    'Access-Control-Allow-Origin': '*',
                    'Access-Control-Allow-Methods': 'GET, POST, PUT, DELETE, OPTIONS',
                    'Access-Control-Allow-Headers': 'Origin, Content-Type, Authorization',
                    'Access-Control-Allow-Credentials': 'true'
                }
            },
            mounted() {
                this.fetchSettings();
            },
            watch: {
                'settings.proxy_on': function (newVal, oldVal) {
                    if (newVal !== oldVal) {
                        this.saveSettings();
                    }
                }
            },
            methods: {
                fetchSettings() {
                    axios.get('http://localhost:8099/settings')
                        .then(response => {
                            this.settings = response.data;
                        })
                        .catch(error => {
                            console.error(error);
                        });
                },
                saveSettings() {
                    axios.post('http://localhost:8099/settings', this.settings)
                        .then(response => {
                            console.log('Settings saved successfully:', response.data);
                        })
                        .catch(error => {
                            console.error('Error while saving settings:', error);
                        });
                },
                addKeyword() {
                    if (this.newKeyword.trim() !== '') {
                        this.settings.key_words.push(this.newKeyword);
                        this.newKeyword = ''; // Reset the input field
                        this.settings.cookie = '';
                        this.saveSettings();
                    }
                },
                removeKeyword(index) {
                    this.settings.key_words.splice(index, 1);
                    this.saveSettings();
                },
                startGetPoints() {
                    axios.post('http://localhost:8099/start', this.settings)
                        .then(response => {
                            console.log('Start successfully:', response.data);
                        })
                        .catch(error => {
                            console.error('Error Start settings:', error);
                        });
                }
            }
        });
    </script>
</body>

</html>`
