webpackJsonp([8],{"+SOK":function(e,l,r){"use strict";Object.defineProperty(l,"__esModule",{value:!0});var a=r("kl55"),t=r.n(a);for(var i in a)"default"!==i&&function(e){r.d(l,e,function(){return a[e]})}(i);var o=r("MrHA");var s=function(e){r("ANwl")},u=r("VU/8")(t.a,o.a,!1,s,"data-v-6675fe02",null);l.default=u.exports},ANwl:function(e,l){},MrHA:function(e,l,r){"use strict";var a={render:function(){var e=this,l=e.$createElement,r=e._self._c||l;return r("div",{staticClass:"fs3_back"},[e._m(0),e._v(" "),r("div",{staticClass:"fs3_cont"},[r("el-form",{ref:"ruleForm",staticClass:"demo-ruleForm",attrs:{model:e.ruleForm,rules:e.rules}},[r("el-form-item",{attrs:{label:"Backup plan name:",prop:"name"}},[r("el-input",{model:{value:e.ruleForm.name,callback:function(l){e.$set(e.ruleForm,"name",l)},expression:"ruleForm.name"}})],1),e._v(" "),r("el-form-item",{attrs:{label:"Choose your backup frequency:",prop:"frequency"}},[r("el-select",{attrs:{placeholder:""},model:{value:e.ruleForm.frequency,callback:function(l){e.$set(e.ruleForm,"frequency",l)},expression:"ruleForm.frequency"}},e._l(e.ruleForm.frequencyOptions,function(e){return r("el-option",{key:e.value,attrs:{label:e.label,value:e.value}})}),1)],1),e._v(" "),r("el-form-item",{attrs:{label:"Choose your backup region:",prop:"region"}},[r("el-select",{attrs:{placeholder:""},model:{value:e.ruleForm.region,callback:function(l){e.$set(e.ruleForm,"region",l)},expression:"ruleForm.region"}},e._l(e.ruleForm.regionOptions,function(e){return r("el-option",{key:e.value,attrs:{label:e.label,value:e.value}})}),1)],1),e._v(" "),r("el-form-item",{attrs:{label:"Price:",prop:"price"}},[r("el-input",{staticClass:"input",attrs:{onkeyup:"value=value.replace(/^\\D*(\\d*(?:\\.\\d{0,20})?).*$/g, '$1')"},model:{value:e.ruleForm.price,callback:function(l){e.$set(e.ruleForm,"price",l)},expression:"ruleForm.price"}}),e._v(" FIL\n      ")],1),e._v(" "),r("el-form-item",{attrs:{label:"Duration:",prop:"duration"}},[r("el-input",{staticClass:"input",attrs:{onkeyup:"value=value.replace(/^(0+)|[^\\d]+/g,'')"},model:{value:e.ruleForm.duration,callback:function(l){e.$set(e.ruleForm,"duration",l)},expression:"ruleForm.duration"}}),e._v(" Days\n      ")],1),e._v(" "),r("el-form-item",{attrs:{label:"Verified-Deal:",prop:"verified"}},[r("el-radio",{attrs:{label:"1"},model:{value:e.ruleForm.verified,callback:function(l){e.$set(e.ruleForm,"verified",l)},expression:"ruleForm.verified"}},[e._v("Yes")]),e._v(" "),r("el-radio",{attrs:{label:"2"},model:{value:e.ruleForm.verified,callback:function(l){e.$set(e.ruleForm,"verified",l)},expression:"ruleForm.verified"}},[e._v("No")])],1),e._v(" "),r("el-form-item",{attrs:{label:"Fast-Retrival:",prop:"fastRetirval"}},[r("el-radio",{attrs:{label:"1"},model:{value:e.ruleForm.fastRetirval,callback:function(l){e.$set(e.ruleForm,"fastRetirval",l)},expression:"ruleForm.fastRetirval"}},[e._v("Yes")]),e._v(" "),r("el-radio",{attrs:{label:"2"},model:{value:e.ruleForm.fastRetirval,callback:function(l){e.$set(e.ruleForm,"fastRetirval",l)},expression:"ruleForm.fastRetirval"}},[e._v("No")])],1),e._v(" "),r("el-form-item",[r("el-button",{on:{click:function(l){return e.submitForm("ruleForm")}}},[e._v("Create")])],1)],1)],1),e._v(" "),r("el-dialog",{attrs:{title:e.ruleForm.name,"custom-class":"formStyle",visible:e.dialogVisible,width:e.dialogWidth},on:{"update:visible":function(l){e.dialogVisible=l}}},[r("el-card",{staticClass:"box-card"},[r("div",{staticClass:"statusStyle"},[r("div",{staticClass:"list"},[r("span",[e._v("Backup frequency:")]),e._v(" "+e._s(e.ruleForm.frequency))]),e._v(" "),r("div",{staticClass:"list"},[r("span",[e._v("Backup region:")]),e._v(" "+e._s(e.ruleForm.region))]),e._v(" "),r("div",{staticClass:"list"},[r("span",[e._v("Price:")]),e._v(" "+e._s(e.ruleForm.price)+" FIL")]),e._v(" "),r("div",{staticClass:"list"},[r("span",[e._v("Duration:")]),e._v(" "+e._s(e.ruleForm.duration)+" days")]),e._v(" "),r("div",{staticClass:"list"},[r("span",[e._v("Verified deal:")]),e._v(" "+e._s("2"==e.ruleForm.verified?"No":"Yes"))]),e._v(" "),r("div",{staticClass:"list"},[r("span",[e._v("Fast retrieval:")]),e._v(" "+e._s("2"==e.ruleForm.fastRetirval?"No":"Yes"))])])]),e._v(" "),r("div",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[r("el-button",{on:{click:e.confirm}},[e._v("OK")])],1)],1),e._v(" "),r("el-dialog",{attrs:{title:"Backup Plans","custom-class":"formStyle",visible:e.dialogConfirm,width:e.dialogWidth},on:{"update:visible":function(l){e.dialogConfirm=l}}},[r("span",{staticClass:"span"},[e._v("Your backup has created successfully")]),e._v(" "),r("div",{staticClass:"dialog-footer",attrs:{slot:"footer"},slot:"footer"},[r("el-button",[e._v("VIEW")]),e._v(" "),r("el-button",{on:{click:function(l){e.dialogConfirm=!1}}},[e._v("OK")])],1)])],1)},staticRenderFns:[function(){var e=this.$createElement,l=this._self._c||e;return l("div",{staticClass:"fs3_head"},[l("div",{staticClass:"fs3_head_text"},[l("div",{staticClass:"titleBg"},[this._v("Backup Plans")]),this._v(" "),l("h1",[this._v("Backup Plans")])]),this._v(" "),l("img",{staticClass:"bg",attrs:{src:r("3Msz"),alt:""}})])}]};l.a=a},kl55:function(e,l,r){"use strict";Object.defineProperty(l,"__esModule",{value:!0});var a,t=r("mtWM");(a=t)&&a.__esModule;l.default={data:function(){return{width:document.body.clientWidth>600?"400px":"95%",dialogWidth:document.body.clientWidth<=600?"95%":"50%",dialogVisible:!1,dialogConfirm:!1,ruleForm:{name:"",price:"",duration:"",verified:"2",fastRetirval:"1",frequency:"Backup Daily",frequencyOptions:[{value:"Backup Daily",label:"Backup Daily"},{value:"Backup Weekly",label:"Backup Weekly"}],region:"Global",regionOptions:[{value:"Global",label:"Global"},{value:"Asia",label:"Asia"},{value:"Africa",label:"Africa"},{value:"North America",label:"North America"},{value:"Sorth America",label:"Sorth America"},{value:"Europe",label:"Europe"},{value:"Oceania",label:"Oceania"}]},rules:{name:[{required:!0,message:"Please enter Backup plan Name",trigger:"blur"}],price:[{required:!0,message:"Please enter Price",trigger:"blur"}],duration:[{required:!0,message:"Please enter Duration",trigger:"blur"}]}}},watch:{},methods:{confirm:function(){this.dialogVisible=!1,this.dialogConfirm=!0},submitForm:function(e){var l=this;l.$refs[e].validate(function(e){if(!e)return console.log("error submit!!"),!1;l.dialogVisible=!0})}},mounted:function(){}}}});