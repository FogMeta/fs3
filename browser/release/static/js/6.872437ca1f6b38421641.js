webpackJsonp([6],{"6kc2":function(t,a,e){"use strict";Object.defineProperty(a,"__esModule",{value:!0});var s=n(e("mtWM")),i=n(e("PJh5"));function n(t){return t&&t.__esModule?t:{default:t}}a.default={data:function(){return{dialogWidth:document.body.clientWidth<=600?"95%":"600px",dialogIndex:NaN,dialogVisible:!1,width:document.body.clientWidth>600?"400px":"95%",ruleForm:{},plan_list:[],loading:!1,parma:{limit:10,offset:1,total:0},bodyWidth:document.documentElement.clientWidth<1024,deleteConfirm:!1}},watch:{},methods:{planSubmit:function(t,a){this.dialogIndex=t,this.dialogVisible=!0,this.ruleForm=a},handleClose:function(t){this.dialogIndex=NaN,this.dialogVisible=!1},deleteAct:function(t){this.deleteConfirm=!0,this.ruleForm=t},planStatus:function(t,a){var e=this;e.deleteConfirm=!1,e.dialogVisible=!1,e.dialogIndex=NaN,e.loading=!0;var n={BackupPlanId:t.ID,Status:a?"Deleted":"Enabled"==t.Status?"Disabled":"Enabled"};s.default.post(e.data_api+"/minio/backup/update/plan",n,{headers:{Authorization:"Bearer "+e.$store.getters.accessToken}}).then(function(t){var a=t.data;if("success"!=a.status)return e.$message.error(a.message),e.loading=!1,!1;e.ruleForm=a.data,e.ruleForm.CreatedOn=e.ruleForm.CreatedOn?(0,i.default)(new Date(parseInt(e.ruleForm.CreatedOn/1e3))).format("YYYY-MM-DD HH:mm:ss"):"-",e.ruleForm.UpdatedOn=e.ruleForm.UpdatedOn?(0,i.default)(new Date(parseInt(e.ruleForm.UpdatedOn/1e3))).format("YYYY-MM-DD HH:mm:ss"):"-",e.ruleForm.LastBackupOn=e.ruleForm.LastBackupOn?(0,i.default)(new Date(parseInt(e.ruleForm.LastBackupOn/1e3))).format("YYYY-MM-DD HH:mm:ss"):"-",e.getData()}).catch(function(t){console.log(t),e.loading=!1})},sort:function(t){return t.sort(function(t,a){return t.ID-a.ID})},handleCurrentChange:function(t){this.parma.offset=t,this.getData()},getData:function(){var t=this;t.loading=!0;var a=t.data_api+"/minio/backup/retrieve/plan",e={Offset:t.parma.offset>0?t.parma.offset-1:t.parma.offset,Limit:t.parma.limit,Status:["Enabled","Disabled"]};s.default.post(a,e,{headers:{Authorization:"Bearer "+t.$store.getters.accessToken}}).then(function(a){var e=a.data;if("success"!=e.status)return t.loading=!1,t.$message.error(e.message),!1;t.parma.total=e.data.TotalVolumeBackupPlan,t.plan_list=e.data.backupPlans,t.plan_list.map(function(t){t.CreatedOn=t.CreatedOn?(0,i.default)(new Date(parseInt(t.CreatedOn/1e3))).format("YYYY-MM-DD HH:mm:ss"):"-",t.UpdatedOn=t.UpdatedOn?(0,i.default)(new Date(parseInt(t.UpdatedOn/1e3))).format("YYYY-MM-DD HH:mm:ss"):"-",t.LastBackupOn=t.LastBackupOn?(0,i.default)(new Date(parseInt(t.LastBackupOn/1e3))).format("YYYY-MM-DD HH:mm:ss"):"-"}),setTimeout(function(){t.sort(t.plan_list),t.loading=!1},500)}).catch(function(a){t.loading=!1,console.log(a)})}},mounted:function(){this.getData()}}},"K+81":function(t,a,e){"use strict";Object.defineProperty(a,"__esModule",{value:!0});var s=e("6kc2"),i=e.n(s);for(var n in s)"default"!==n&&function(t){e.d(a,t,function(){return s[t]})}(n);var l=e("Us/S");var r=function(t){e("ZzJD")},o=e("VU/8")(i.a,l.a,!1,r,"data-v-0227496a",null);a.default=o.exports},"Us/S":function(t,a,e){"use strict";var s={render:function(){var t=this,a=t.$createElement,e=t._self._c||a;return e("div",{staticClass:"fs3_back"},[t._m(0),t._v(" "),e("div",{staticClass:"fs3_cont"},[e("div",{staticClass:"backupCreate"},[e("router-link",{attrs:{to:{name:"my_account_backupPlans"}}},[t._v("Create Backup Plan")])],1),t._v(" "),t.plan_list.length>0?e("div",{directives:[{name:"loading",rawName:"v-loading",value:t.loading,expression:"loading"}]},t._l(t.plan_list,function(a,s){return e("el-card",{key:s,staticClass:"box-card"},[e("div",{staticClass:"title"},[t._v(t._s(a.Name))]),t._v(" "),e("div",{staticClass:"button"},["Created"==a.Status?e("div",{staticClass:"statusStyle",staticStyle:{color:"#0a318e"}},[t._v("\n              "+t._s(a.Status)+"\n          ")]):"Running"==a.Status?e("div",{staticClass:"statusStyle",staticStyle:{color:"#ffb822"}},[t._v("\n              "+t._s(a.Status)+"\n          ")]):"Enabled"==a.Status?e("div",{staticClass:"statusStyle",staticStyle:{color:"#67c23a"}},[t._v("\n              "+t._s(a.Status)+"\n          ")]):"Completed"==a.Status?e("div",{staticClass:"statusStyle",staticStyle:{color:"#1dc9b7"}},[t._v("\n              "+t._s(a.Status)+"\n          ")]):"Stopped"==a.Status?e("div",{staticClass:"statusStyle",staticStyle:{color:"#f56c6c"}},[t._v("\n              "+t._s(a.Status)+"\n          ")]):"Disabled"==a.Status?e("div",{staticClass:"statusStyle",staticStyle:{color:"#e41f1f"}},[t._v("\n              "+t._s(a.Status)+"\n          ")]):e("div",{staticClass:"statusStyle",staticStyle:{color:"rgb(255, 184, 34)"}},[t._v("\n              "+t._s(a.Status)+"\n          ")]),t._v(" "),e("el-button",{class:{active:t.dialogIndex==s,btnSub:!0},on:{click:function(e){return t.planSubmit(s,a)}}},[t._v("View details")]),t._v(" "),e("el-button",{attrs:{type:"danger"},on:{click:function(e){return t.deleteAct(a,1)}}},[t._v("Delete")])],1)])}),1):e("div",{staticStyle:{"text-align":"center"}},[t._v("No Data")])]),t._v(" "),e("div",{staticClass:"form_pagination"},[e("div",{staticClass:"pagination"},[e("el-pagination",{attrs:{"hide-on-single-page":"",total:t.parma.total,"page-size":t.parma.limit,"current-page":t.parma.offset,"pager-count":t.bodyWidth?5:7,background:"",layout:t.bodyWidth?"prev, pager, next":"total, prev, pager, next, jumper"},on:{"current-change":t.handleCurrentChange}})],1)]),t._v(" "),e("el-dialog",{attrs:{title:t.ruleForm.Name,"custom-class":"formStyle",visible:t.dialogVisible,width:t.dialogWidth,"before-close":t.handleClose},on:{"update:visible":function(a){t.dialogVisible=a}}},[e("el-card",{staticClass:"box-card"},[e("div",{staticClass:"statusStyle"},[e("div",{staticClass:"list"},[e("span",[t._v("Backup Plan ID:")]),t._v(" "+t._s(t.ruleForm.ID))]),t._v(" "),e("div",{staticClass:"list"},[e("span",[t._v("Backup frequency:")]),t._v(" "+t._s("1"==t.ruleForm.Interval?"Backup Daily":"Backup Weekly"))]),t._v(" "),e("div",{staticClass:"list"},[e("span",[t._v("Price:")]),t._v(" "+t._s(t.ruleForm.Price)+" FIL")]),t._v(" "),e("div",{staticClass:"list"},[e("span",[t._v("Duration:")]),t._v(" "+t._s(t.ruleForm.Duration/24/60/2)+" days")]),t._v(" "),e("div",{staticClass:"list"},[e("span",[t._v("Verified deal:")]),t._v(" "+t._s(t.ruleForm.VerifiedDeal?"Yes":"No"))]),t._v(" "),e("div",{staticClass:"list"},[e("span",[t._v("Fast retrieval:")]),t._v(" "+t._s(t.ruleForm.FastRetrieval?"Yes":"No"))]),t._v(" "),e("div",{staticClass:"list"},[e("span",[t._v("Create Date:")]),t._v(" "+t._s(t.ruleForm.CreatedOn))]),t._v(" "),e("div",{staticClass:"list"},[e("span",[t._v("Last Update:")]),t._v(" "+t._s(t.ruleForm.UpdatedOn))]),t._v(" "),e("div",{staticClass:"list"},[e("span",[t._v("Last Backup Date:")]),t._v(" "+t._s(t.ruleForm.LastBackupOn))])])]),t._v(" "),e("div",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[e("el-button",{attrs:{type:t.ruleForm.Status&&"enabled"==t.ruleForm.Status.toLowerCase()?"danger":"info",disabled:!t.ruleForm.Status||"enabled"!=t.ruleForm.Status.toLowerCase()},on:{click:function(a){return t.planStatus(t.ruleForm)}}},[t._v("STOP")]),t._v(" "),e("el-button",{attrs:{type:t.ruleForm.Status&&"enabled"==t.ruleForm.Status.toLowerCase()?"info":"success",disabled:!(!t.ruleForm.Status||"enabled"!=t.ruleForm.Status.toLowerCase())},on:{click:function(a){return t.planStatus(t.ruleForm)}}},[t._v("START")])],1)],1),t._v(" "),e("el-dialog",{attrs:{title:"Tips","custom-class":"formStyle",visible:t.deleteConfirm,width:t.dialogWidth},on:{"update:visible":function(a){t.deleteConfirm=a}}},[e("span",{staticClass:"span"},[t._v("Are you sure you want to delete?")]),t._v(" "),e("div",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[e("el-button",{attrs:{type:"danger"},on:{click:function(a){return t.planStatus(t.ruleForm,1)}}},[t._v("Yes, delete")]),t._v(" "),e("el-button",{attrs:{type:"info"},on:{click:function(a){t.deleteConfirm=!1}}},[t._v("No, close the window")])],1)])],1)},staticRenderFns:[function(){var t=this.$createElement,a=this._self._c||t;return a("div",{staticClass:"fs3_head"},[a("div",{staticClass:"fs3_head_text"},[a("div",{staticClass:"titleBg"},[this._v("My Plans")]),this._v(" "),a("h1",[this._v("My Plans")])]),this._v(" "),a("img",{staticClass:"bg",attrs:{src:e("3Msz"),alt:""}})])}]};a.a=s},ZzJD:function(t,a){}});