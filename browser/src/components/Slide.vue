<template>
    <div class="slide" :class="{'slide_stretch': menuStretch}" @click="caozuoFun()">
        <div class="slideScroll" v-if="!menuStretch">
            <div class="fes-header">
                <img :src="logo" alt="">
                <svg t="1639358873000" @click="menuToggleStreth(true)" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="2874" width="200" height="200"><path d="M470.784 690.858667H86.613333c-21.418667 0-38.741333-19.114667-38.741333-42.666667s17.322667-42.666667 38.741333-42.666667h384.085334c21.418667 0 38.741333 19.114667 38.741333 42.666667s-17.322667 42.666667-38.656 42.666667zM470.784 449.877333H86.613333c-21.418667 0-38.741333-19.114667-38.741333-42.666666s17.322667-42.666667 38.741333-42.666667h384.085334c21.418667 0 38.741333 19.114667 38.741333 42.666667s-17.322667 42.666667-38.656 42.666666zM934.058667 208.981333H84.394667c-21.418667 0-38.741333-19.114667-38.741334-42.666666s17.322667-42.666667 38.741334-42.666667h849.664c21.418667 0 38.741333 19.114667 38.741333 42.666667s-17.322667 42.666667-38.741333 42.666666zM931.754667 931.754667H82.090667c-21.418667 0-38.741333-19.114667-38.741334-42.666667s17.322667-42.666667 38.741334-42.666667h849.664c21.418667 0 38.741333 19.114667 38.741333 42.666667s-17.322667 42.666667-38.741333 42.666667zM931.754667 724.992c-5.546667 0-11.093333-1.28-16.298667-3.925333L610.218667 565.418667c-13.738667-6.997333-22.442667-22.101333-22.442667-38.741334 0-16.64 8.789333-31.744 22.442667-38.741333l305.237333-155.648c12.032-6.144 26.026667-5.12 37.12 2.730667 11.178667 7.850667 17.92 21.418667 17.92 35.925333v311.296c0 14.592-6.741333 28.074667-17.92 35.925333-6.229333 4.522667-13.482667 6.826667-20.821333 6.826667z" p-id="2875" fill="#ffffff"></path></svg>
            </div>
            <div class="fes-search">
                <el-input
                    placeholder="Search Buckets..."
                    prefix-icon="el-icon-search"
                    v-model="search"
                    @input="searchBucketFun">
                </el-input>
                <el-row>
                    <el-col :span="24" v-for="(item, index) in minioListBucketsAll.buckets" :key="index" :class="{'active': item.name == currentBucket && allActive}" @click.native="getListBucket(item.name, true)">
                        <div>
                            <i class="iconfont icon-harddriveyingpan"></i>
                            {{item.name}}
                        </div>
                        <i class="caozuo el-icon-more" @click.stop="caozuoFun(index, item.name)"></i>

                        <ul v-if="item.show && homeClick">
                            <li @click.stop="dialogFun(item.name, index)">Edit policy</li>
                            <li @click="backupFun">Backup to Filecoin</li>
                            <li @click="retrievalFun">Retrieval</li>
                            <li @click.stop="dialogDeleteFun(item.name, index)">Delete</li>
                        </ul>
                    </el-col>

                    <!-- <el-col :span="24" class="active"
                    style="margin-top:0.2rem;justify-content: center;padding: 0.1rem 0;color: #fff" @click.native="getListBucket('', false, false, true)">
                    All Deals
                    </el-col> -->
                </el-row>
            </div>
            <div class="fs3_backup">
                <div class="introduce">
                    <router-link :to="{name: 'backup'}" :style="{'color': introduceColor?'#2f85e5':'#fff'}">FS3 Backup</router-link>
                </div>
                <div class="introRouter">
                    <router-link :to="{name: 'my_account_dashboard'}" :class="{'active': activeTree == '2'}">Dashboard</router-link>
                    <router-link :to="{name: 'my_account_myPlans'}" :class="{'active': activeTree == '3'}">Backup Plans</router-link>
                    <router-link :to="{name: 'my_account_jobs'}" :class="{'active': activeTree == '4'}">Jobs</router-link>
                </div>
                <!-- :default-checked-keys="activeTree" -->
                <!-- <el-tree :data="dataBackup" :props="defaultProps" @node-click="handleNodeClick"
                    node-key="id" ref="my-tree" default-expand-all
                    :default-expanded-keys="activeTree?[1]:[]"
                    :current-node-key="activeTree"></el-tree> -->

            </div>
        </div>
        <div class="slide_stretch_show" v-if="menuStretch">
            <img src="@/assets/images/logo_small.png" alt="">
            <svg t="1639358873000" @click="menuToggleStreth(false)" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg" p-id="2874" width="200" height="200"><path d="M470.784 690.858667H86.613333c-21.418667 0-38.741333-19.114667-38.741333-42.666667s17.322667-42.666667 38.741333-42.666667h384.085334c21.418667 0 38.741333 19.114667 38.741333 42.666667s-17.322667 42.666667-38.656 42.666667zM470.784 449.877333H86.613333c-21.418667 0-38.741333-19.114667-38.741333-42.666666s17.322667-42.666667 38.741333-42.666667h384.085334c21.418667 0 38.741333 19.114667 38.741333 42.666667s-17.322667 42.666667-38.656 42.666666zM934.058667 208.981333H84.394667c-21.418667 0-38.741333-19.114667-38.741334-42.666666s17.322667-42.666667 38.741334-42.666667h849.664c21.418667 0 38.741333 19.114667 38.741333 42.666667s-17.322667 42.666667-38.741333 42.666666zM931.754667 931.754667H82.090667c-21.418667 0-38.741333-19.114667-38.741334-42.666667s17.322667-42.666667 38.741334-42.666667h849.664c21.418667 0 38.741333 19.114667 38.741333 42.666667s-17.322667 42.666667-38.741333 42.666667zM931.754667 724.992c-5.546667 0-11.093333-1.28-16.298667-3.925333L610.218667 565.418667c-13.738667-6.997333-22.442667-22.101333-22.442667-38.741334 0-16.64 8.789333-31.744 22.442667-38.741333l305.237333-155.648c12.032-6.144 26.026667-5.12 37.12 2.730667 11.178667 7.850667 17.92 21.418667 17.92 35.925333v311.296c0 14.592-6.741333 28.074667-17.92 35.925333-6.229333 4.522667-13.482667 6.826667-20.821333 6.826667z" p-id="2875" fill="#ffffff"></path></svg>
            <!-- <i class="iconfont icon-ziyuan" @click="menuToggleStreth(false)"></i> -->
        </div>
        <div class="fes-host">
            <div class="fesHostLink" v-if="!menuStretch">
                <i class="iconfont icon-diqiu"></i>
                <a href="/">{{location}}</a>
            </div>
            <div class="fesHostLogout" @click="logout">
                <router-link to="/fs3/login" id="logout"><i class="iconfont icon-signout"></i></router-link>
            </div>
        </div>

        <el-dialog :title="titlePolicy" :visible.sync="dialogFormVisible" custom-class="policyStyle">
            <el-form ref="dynamicValidateForm" class="demo-dynamic">
                <el-form-item>
                    <el-input v-model="dynamicValidateForm.value" placeholder="Prefix"></el-input>
                    <el-select v-model="dynamicValidateForm.valueSelect" placeholder="">
                        <el-option
                        v-for="item in dynamicValidateForm.options"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value">
                        </el-option>
                    </el-select>
                    <el-button type="primary" @click="addPolicies">Add</el-button>
                </el-form-item>
                <el-form-item v-for="(domain, index) in bucketPolicies.policies" :key="index">
                    <!--el-input v-model="domain.prefix" placeholder="Prefix" disabled></el-input-->
                    <div class="el-input">
                      <div class="el-input__inner">{{domain.prefix?domain.prefix:'*'}}</div>
                    </div>
                    <el-select v-model="domain.policy" placeholder="" disabled>
                        <el-option
                        v-for="item in dynamicValidateForm.options"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value">
                        </el-option>
                    </el-select>
                    <el-button type="danger" @click.prevent="removePolicies(domain)">Remove</el-button>
                </el-form-item>
            </el-form>
        </el-dialog>


    </div>
</template>
<script>
import axios from 'axios'
export default {
    data() {
        return {
            postUrl: this.data_api + `/minio/webrpc`,
            logo: require("@/assets/images/logo.png"),
            activeIndex: '1',
            mobileMenuShow: false,
            search: '',
            location: window.location.host,
            minioListBucketsAll: {},
            dialogFormVisible: false,
            titlePolicy: 'Bucket Policy ',
            dynamicValidateForm: {
                value: '',
                options: [{
                    value: 'readonly',
                    label: 'Read Only'
                },{
                    value: 'writeonly',
                    label: 'Write Only'
                },{
                    value: 'readwrite',
                    label: 'Read and Write'
                }],
                valueSelect: 'readonly'
            },
            bucketPolicies: {
                policies: []
            },
            shareDialog: false,
            shareObjectShow: true,
            shareFileShow: false,
            sendApi: 1,
            retrievalDialog: false,
            allActive: true,
            dataBackup: [{
                label: 'My account',
                id: 1,
                children: [{
                    label: 'Dashboard',
                    id: 2,
                    urlName: 'my_account_dashboard'
                },{
                    label: 'Backup Plans',
                    id: 3,
                    urlName: 'my_account_backupPlans'
                },{
                    label: 'My Plans',
                    id: 4,
                    urlName: 'my_account_myPlans'
                }]
            }],
            activeTree: '',
            defaultProps: {
                children: 'children',
                label: 'label'
            },
            introduceColor: false,
            menuStretch: false
        };
    },
    props: ['minioListBuckets', 'currentBucket', 'homeClick'],
    components: {},
    computed: {
        email() {
            return this.$store.state.user.email
        },
    },
    watch: {
        $route: function (to, from) {
            this.productName()
            if(this.bodyWidth){
                this.collapse = true
                this.collapseChage();
            }
        },
        'minioListBuckets': function (to, from) {
            this.getMinioData()
        },
        // activeTree(id) {
        //     // Tree 内部使用了 Node 类型的对象来包装用户传入的数据，用来保存目前节点的状态。可以用 $refs 获取 Tree 实例
        //     if (id.toString()) {
        //         this.$refs["my-tree"].setCurrentKey(id);
        //     } else {
        //         this.$refs["my-tree"].setCurrentKey(null);
        //     }
        // }
    },
    methods: {
      menuToggleStreth(stretch) {
        this.menuStretch = stretch
        this.$emit('getMenuStretch', stretch);
      },
      logout() {
        var _this = this;

        let dataGetDiscoveryDoc = {
            id: 1,
            jsonrpc: "2.0",
            method: "web.GetDiscoveryDoc",
            params:{}
        }
        axios.post(_this.postUrl, dataGetDiscoveryDoc, {headers: {
            'Authorization':"Bearer "+ _this.$store.getters.accessToken
        }}).then((response) => {
            let json = response.data
            let error = json.error
            let result = json.result
            if (error) {
                _this.$message.error(error.message);
                return false
            }

            _this.$store.dispatch("FedLogOut").then(() => {
                _this.$router.replace({ name: 'login' })
            });

        }).catch(function (error) {
            console.log(error);
        });
      },
      handleNodeClick(data) {
        if(data.urlName){
            this.$router.push({name: data.urlName})
            this.getListBucket('', false, true)
        } 
      },
      productName() {
        let _this = this
        _this.introduceColor = _this.$route.name == 'backup'?true:false
        _this.activeTree = ''
        if(_this.$route.name.indexOf('my_account') > -1){
            if(_this.$route.name == 'my_account_backupPlans' || _this.$route.name == 'my_account_myPlans') {
                _this.activeTree = '3'
            }else if(_this.$route.name == 'my_account_jobs') {
                _this.activeTree = '4'
            }else {
                _this.activeTree = '2'
            }
        }
      },
      getshareDialog(shareDialog) {
        this.shareDialog = shareDialog
      },
      getretrievalDialog(retrievalDialog) {
        this.retrievalDialog = retrievalDialog
      },
      backupFun() {
        this.shareDialog = true
        this.shareObjectShow = false
        this.shareFileShow = true
        this.$emit('getshareHome', true, false, true);
      },
      retrievalFun(){
        this.$emit('getretrievalHome', true);
      },
      removePolicies(content) {
        let _this = this
        let dataSetBucketPolicy = {
                id: 1,
                jsonrpc: "2.0",
                method: "web.SetBucketPolicy",
                params:{
                    bucketName: content.bucket,
                    policy: "none",
                    prefix: content.prefix
                }
            }
            axios.post(_this.postUrl, dataSetBucketPolicy, {headers: {
                'Authorization':"Bearer "+ _this.$store.getters.accessToken
            }}).then((response) => {
                let json = response.data
                let error = json.error
                let result = json.result
                if (error) {
                    _this.$message.error(error.message);
                    return false
                }
                _this.dynamicValidateForm.value = ''
                _this.getListAllBucketPolicies(_this.currentBucket)

            }).catch(function (error) {
                console.log(error);
            });
      },
      addPolicies() {
            let _this = this
            let $hgh
            if(_this.bucketPolicies.policies) {
              _this.bucketPolicies.policies.map(item => {
                  if((item.prefix == _this.dynamicValidateForm.value && item.policy == _this.dynamicValidateForm.valueSelect) || ((!item.prefix) && _this.dynamicValidateForm.value == '*') && item.policy !== _this.dynamicValidateForm.valueSelect){
                      _this.$message({
                          message: 'Policy for this prefix already exists.',
                          type: 'warning',
                          showClose: true
                      });
                      $hgh = true
                      return false
                  }else{
                      // console.log(item.prefix, _this.dynamicValidateForm.value, item.policy, _this.dynamicValidateForm.valueSelect);
                  }
              })
           }

           if(!$hgh){
            _this.setPolicyChange()
           }

      },
      setPolicyChange() {
            let _this = this
            let dataSetBucketPolicy = {
                id: 1,
                jsonrpc: "2.0",
                method: "web.SetBucketPolicy",
                params:{
                    bucketName: _this.currentBucket,
                    policy: _this.dynamicValidateForm.valueSelect,
                    prefix: _this.dynamicValidateForm.value
                }
            }
            axios.post(_this.postUrl, dataSetBucketPolicy, {headers: {
                'Authorization':"Bearer "+ _this.$store.getters.accessToken
            }}).then((response) => {
                let json = response.data
                let error = json.error
                let result = json.result
                if (error) {
                    _this.$message.error(error.message);
                    return false
                }
                _this.dynamicValidateForm.value = ''
                _this.getListAllBucketPolicies(_this.currentBucket)

            }).catch(function (error) {
                console.log(error);
            });
      },
      dialogFun(name, index) {
        let _this = this
        _this.titlePolicy = 'Bucket Policy (' + name + ')'
        _this.dialogFormVisible = true
        if(_this.minioListBucketsAll) {
          _this.minioListBucketsAll.buckets.map((item, i) => {
              item.show = false;
          })
        }

        _this.getListAllBucketPolicies(name)
      },
      dialogDeleteFun(name, index) {
        let _this = this
        let dataDeleteBucket = {
            id: 1,
            jsonrpc: "2.0",
            method: "web.DeleteBucket",
            params:{
                bucketName: name,
            }
        }
        axios.post(_this.postUrl, dataDeleteBucket, {headers: {
            'Authorization':"Bearer "+ _this.$store.getters.accessToken
        }}).then((response) => {
            let json = response.data
            let error = json.error
            let result = json.result
            if (error) {
                _this.$message.error(error.message);
                return false
            }

            _this.$emit('getListBuckets');

        }).catch(function (error) {
            console.log(error);
        });
      },
      getListAllBucketPolicies(name) {
        let _this = this
        let dataListAllBucketPolicies = {
            id: 1,
            jsonrpc: "2.0",
            method: "web.ListAllBucketPolicies",
            params:{
                bucketName: name,
            }
        }
        axios.post(_this.postUrl, dataListAllBucketPolicies, {headers: {
            'Authorization':"Bearer "+ _this.$store.getters.accessToken
        }}).then((response) => {
            let json = response.data
            let error = json.error
            let result = json.result
            if (error) {
                _this.$message.error(error.message);
                return false
            }
            _this.bucketPolicies = result

        }).catch(function (error) {
            console.log(error);
        });
      },
      handleSelect(key, keyPath) {
        //console.log(key, keyPath);
      },
      mobileMenuFun(){
        let _this = this;
        _this.mobileMenuShow=!_this.mobileMenuShow;
        if(_this.mobileMenuShow){
          document.body.style.height = '100vh'
          document.body.style['overflow-y'] = 'hidden'
        }else{
          document.body.style.height = 'auto'
          document.body.style['overflow-y'] = 'auto'
        }
      },
      caozuoFun(index, name) {
        let _this = this;
        _this.$nextTick(() => {
          if(_this.minioListBucketsAll.buckets) {
            _this.minioListBucketsAll.buckets.map((item, i) => {
                item.show = false;
                if(i == index){
                    item.show = true
                    _this.$emit('homeClickFun', true)
                }
            })
           }
        })
        if(name){
            _this.getListBucket(name, true)
        }
      },
      searchBucketFun() {
          let _this = this
          if(_this.search){
            _this.minioListBucketsAll.buckets = []
            if(_this.minioListBuckets.buckets) {
              _this.minioListBuckets.buckets.map(item => {
                  if(item.name.indexOf(_this.search) >= 0){
                      _this.minioListBucketsAll.buckets.push(item)
                  }
              })
            }
          }else{
              _this.minioListBucketsAll = JSON.parse(JSON.stringify(_this.minioListBuckets))
          }
      },
      getListBucket(name, allDeal, silde, push) {
          this.$emit('getminioListBucket', name, allDeal, silde, push);
          this.allActive = allDeal ? true : false
      },
      getMinioData() {
        let _this = this;
        if(_this.minioListBuckets && _this.minioListBuckets.buckets){
            _this.minioListBuckets.buckets.map(item => {
                item.show = false;
            })
            _this.minioListBucketsAll = JSON.parse(JSON.stringify(_this.minioListBuckets))

        }else{
          _this.minioListBucketsAll = JSON.parse(JSON.stringify(_this.minioListBuckets))
        }
      }
    },
    mounted() {
      this.getMinioData()
      this.productName()
    },
};
</script>
<style lang="scss" scoped>
.slide{
    position: relative;
    width: 3.2rem;
    background-color: #003040;
    height: 100%;
    overflow: hidden;
    padding: 0;
    transition: all;
    transition-duration: .3s;
    .slideScroll{
        position: relative;
        height: calc(100% - 0.6rem);
        overflow: hidden;
        overflow-y: scroll;
        scrollbar-color: #ccc #002a39;
        scrollbar-width: 4px;
        scrollbar-track-color: transparent;
        -ms-scrollbar-track-color: transparent;
        &::-webkit-scrollbar-track {
            background: #003040;
        }
        &::-webkit-scrollbar {
            width: 4px;
            background: #002a39;
        }
        &::-webkit-scrollbar-thumb {
            background: #ccc;
        }
    }
    .fes-header{
        display: flex;
        justify-content: space-between;
        align-content: center;
        width: calc(100% - 0.4rem);
        padding: 0.25rem 0.2rem;
        img{
            width: auto;
            max-width: 100%;
            height: 0.35rem;
        }
        svg{
            font-size: 0.25rem;
            width: 0.25rem;
            height: 0.25rem;
            margin: 0.05rem 0 0;
            color: rgba(255,255,255,0.75);
            cursor: pointer;
        }
        h2{
            margin: 10px 0 0 13px;
            font-weight: 400;
            color: #fff;
            font-size: 0.2rem;
        }
    }
    .fs3_backup{
        margin: 0.1rem 0 0;
        .introduce{
            margin: 0 0 0.05rem;
            text-indent: 0.2rem;
            background: #002a39;
            // font-family: 'm-semibold';
            font-weight: bold;
            a{
                display: block;
                line-height: 3;
                font-size: 0.14rem;
                color: #2f85e5;
                @media screen and (max-width:999px){
                  font-size: 13px;
                  line-height: 3.5;
                }
            }
        }
        .introRouter{
            font-size: 0.14rem;
            @media screen and (max-width:999px){
                font-size: 13px;
            }
            a{
                display: block;
                padding: 0.07rem 0.2rem;
                color: rgba(255, 255, 255, 0.85);
                font-size: inherit;
                @media screen and (max-width:999px){
                    padding: 8px 0.2rem;
                }
                &:hover{
                    color: #7ecef4;
                    background-color: rgba(0,0,0,.1);
                }
            }
            .active{
                color: #7ecef4;
            }
        }
        .el-tree /deep/{
            padding: 0 0.25rem;
            background: transparent;
            color: #fff;
            .el-tree-node {
                .el-tree-node__content{
                    height: auto;
                    background: transparent !important;
                    margin: 0 0 0.08rem;
                    .el-tree-node__expand-icon{
                        padding: 0 0.05rem;
                        &:before{
                            font-size: 0.2rem;
                        }
                    }
                    .el-tree-node__label{
                        font-size: 0.15rem;
                        @media screen and (max-width:999px){
                            font-size: 14px;
                        }
                    }
                    &:hover{
                        color: #5f9dcc;
                    }
                }
                .el-tree-node__children{
                    .el-tree-node__content{
                        .el-tree-node__label{
                            font-size: 0.14rem;
                            line-height: 1.5;
                            @media screen and (max-width:999px){
                                font-size: 13px;
                            }
                        }
                    }
                }
                .is-current, .is-checked{
                        color: #5f9dcc;
                }
            }
        }
    }
    .fes-search{
        // height: calc(100% - 1.7rem);
        .el-input /deep/{
            display: block;
            width: calc(100% - 0.4rem);
            margin: 0 0.2rem;
            clear: both;
            .el-input__inner{
                padding-left: 25px;
                background-color: transparent;
                box-shadow: none;
                border: 0;
                border-radius: 0;
                border-bottom: 1px solid rgba(255, 255, 255, 0.1);
                color: #fff;
                text-align: left;
                font-size: 0.14rem;
                font-family: inherit;
                @media screen and (max-width:999px){
                    font-size: 13px;
                }
            }
            .el-input__prefix{
                color: #fff;
                left: 0;
                .el-input__icon{
                    text-align: left;
                }
            }
        }
        .el-row /deep/{
            margin-top: 0.2rem;
            margin-left: 0;
            margin-right: 0;
            font-size: 0.13rem;
            // height: calc(100% - 1.3rem);
            // overflow: hidden;
            // overflow-y: scroll;
            .el-col{
                position: relative;
                display: flex;
                align-items: center;
                justify-content: space-between;
                padding: 0.08rem 0.05rem 0.08rem 0.2rem;
                color: rgba(255, 255, 255, 0.75);
                word-wrap: break-word;
                font-size: 0.14rem;
                cursor: pointer;
                @media screen and (max-width:999px){
                    padding: 8px 0.05rem 8px 0.2rem;
                    font-size: 13px;
                }
                div{
                    display: flex;
                    align-items: center;
                }
                i{
                    font-size: 0.18rem;
                    margin-right: 0.08rem;
                    color: rgba(255, 255, 255, 0.85);
                    @media screen and (max-width:999px){
                        font-size: 16px;
                    }
                }
                .caozuo{
                    opacity: 0;
                    float: right;
                    transform: rotate(90deg);
                    font-size: 0.15rem;
                        @media screen and (max-width:999px){
                            font-size: 14px;
                        }
                    color: rgba(255, 255, 255, 0.75);
                    &:hover{
                        color: #fff;
                    }
                }
                ul{
                    position: absolute;
                    right: 0;
                    top: 0;
                    padding: 0.1rem 0;
                    background-color: #fff;
                    border-radius: 0.05rem;
                    z-index: 1000;
                    min-width: 160px;
                    margin: 2px 0 0;
                    list-style: none;
                    font-size: 0.15rem;
                    text-align: left;
                    border: 1px solid transparent;
                    border-radius: 4px;
                    box-shadow: 0 6px 12px rgba(0,0,0,.18);
                    background-clip: padding-box;
                        @media screen and (max-width:999px){
                            font-size: 14px;
                        }
                    li{
                        display: block;
                        padding: 0.08rem 0.1rem;
                        clear: both;
                        font-weight: 400;
                        line-height: 1.42857143;
                        color: #8e8e8e;
                        white-space: nowrap;
                        text-align: right;
                        &:hover{
                            text-decoration: none;
                            color: #333;
                            background-color: rgba(0,0,0,.05);
                        }
                    }
                }
                &:hover{
                    background: rgba(0,0,0,.1);
                    .caozuo{
                        opacity: 1;
                    }
                }
            }
            .active{
                background: rgba(0,0,0,.1);
                font-size: 0.15rem;
                color: #fff;
                @media screen and (max-width:999px){
                    font-size: 13px;
                }
                i{
                    color: #fff;
                }
            }
        }
        .el-row{
          &::-webkit-scrollbar{
              width: 1px;
              height: 1px;
              background-color: #F5F5F5;
          }

          &::-webkit-scrollbar-track {
              box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
              -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
              border-radius: 10px;
              background-color: #F5F5F5;
          }

          &::-webkit-scrollbar-thumb{
              border-radius: 10px;
              box-shadow: inset 0 0 6px rgba(0, 0, 0, .1);
              -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, .1);
              background-color: #c8c8c8;
          }
        }
    }
    .fes-host {
        display: flex;
        justify-content: space-between;
        position: absolute;
        left: 0;
        bottom: 0;
        z-index: 21;
        background-color: rgba(0,0,0,.1);
        font-size: 15px;
        font-weight: 400;
        // width: calc(3.2rem - 0.4rem);
        width: calc(100% - 0.4rem);
        padding: 0.2rem;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
        color: hsla(0,0%,100%,.75);
        transition: all;
        transition-duration: .3s;
        i{
            float: left;
            display: flex;
            justify-content: center;
            align-items: center;
            font-size: 0.2rem;
            margin-right: 10px;
            width: 20px;
            height: 20px;
            color: #888b83;
            // background: url(../assets/images/icon_01.jpg) no-repeat center;
            // background-size: 100%;
            @media screen and (max-width: 999px) {
                font-size: 18px;
            }
        }
        a{
            color: hsla(0,0%,100%,.75);
            font-size: 15px;
            font-weight: 400;
        }
        .fesHostLogout{
            cursor: pointer;
            i{
                margin: 0;
                color: #fff;
                &:hover{
                    color: #7ecef4;
                }
            }
        }
    }
    .el-dialog__wrapper /deep/{
        .policyStyle{
            width: 90%;
            max-width: 600px;
            .el-dialog__header{
                display: flex;
                .el-dialog__title{
                    font-size: 0.15rem;
                    color: #333;
                }
            }
            .el-dialog__body{
                padding: 0.3rem 0;
                .el-form-item{
                    padding: 0.1rem 0.3rem;
                    margin-bottom: 0.1rem;
                    &:nth-child(2n+2){
                        background-color: #f7f7f7;
                    }
                }
                .el-form-item__content{
                    display: flex;
                    line-height: 0.3rem;
                    .el-input{
                        .el-input__inner{
                            height: 0.3rem;
                            padding-left: 0;
                            line-height: 0.3rem;
                            border: 0;
                            border-bottom: 1px solid #f7f7f7;
                            background-color: transparent;
                            border-radius: 0;
                            font-size: 0.13rem;
                            color: #32393f;
                            text-align: left;
                        }
                        .el-input__icon{
                            display: flex;
                            align-items: center;
                        }
                    }
                    .el-input.is-disabled{
                        .el-input__inner{
                            background-color: transparent;
                        }
                    }
                    .el-select{
                        margin: 0 5%;
                    }
                    .el-button{
                        width: 130px;
                        height: 0.3rem;
                        padding: 0;
                        line-height: 0.3rem;
                        color: #fff;
                        font-size: 12px;
                        font-family: inherit;
                        border: 0;
                        border-radius: 0.02rem;
                        text-align: center;
                    }
                }
            }
        }
    }
}
.sliMobile{
    transform: translate3d(0,0,0) !important;
    width: 80%;
    max-width: 400px;
    .fes-header{
        padding: 0;
        height: 65px;
    }
}
.slide_stretch{
    width: 0.65rem;
    .slide_stretch_show{
        display: block;
        padding: 0.2rem 0.1rem;
        img{
            display: block;
            width: 0.4rem;
            max-width: 100%;
            margin: 0 auto 0.15rem;
        }
        svg, i{
            width: 25px;
            height: 25px;
            margin: auto;
            font-size: 25px;
            font-weight: 900;
            color: rgba(255, 255, 255, 0.75);
            display: block;
            text-align: center;
            cursor: pointer;
            transform: rotate(180deg);
            &:hover{
                color: #fff;
            }
        }
    }
}

@media screen and (max-width:1024px){

}
@media screen and (max-width: 999px){
.slide{
    position: fixed;
    left: 0;
    top: 0;
    z-index: 9998;
    transform: translate3d(-3.2rem,0,0);
    .fes-search {
      .el-row /deep/{
        .el-col{
          .caozuo{
            opacity: 1;
          }
        }
      }
    }
}
}
</style>
