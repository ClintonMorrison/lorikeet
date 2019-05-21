(window.webpackJsonp=window.webpackJsonp||[]).push([[0],{40:function(e,t,a){e.exports=a(91)},45:function(e,t,a){},46:function(e,t,a){},55:function(e,t,a){},57:function(e,t,a){},58:function(e,t,a){},59:function(e,t,a){},60:function(e,t,a){},61:function(e,t,a){},62:function(e,t,a){},65:function(e,t,a){},90:function(e,t,a){},91:function(e,t,a){"use strict";a.r(t);var r=a(0),n=a.n(r),s=a(34),o=a.n(s);a(45);function c(){return n.a.createElement("footer",{className:"cp-footer page-footer"},n.a.createElement("div",{className:"container"},n.a.createElement("div",{className:"row"},n.a.createElement("div",{className:"col s12"},n.a.createElement("a",{className:"grey-text text-lighten-4 right",href:"https://clintonmorrison.com"},"Created by Clinton Morrison.")))))}var i=a(1),l=a(2),u=a(4),m=a(3),d=a(5),p=a(6),h=a(14),v=(a(46),[n.a.createElement("li",{key:"register"},n.a.createElement(p.b,{to:"/register"},"Register")),n.a.createElement("li",{key:"login"},n.a.createElement(p.b,{to:"/login"},"Login"))]),f=[n.a.createElement("li",{key:"account"},n.a.createElement(p.b,{to:"/account"},"My Account")),n.a.createElement("li",{key:"passwords"},n.a.createElement(p.b,{to:"/passwords"},"My Passwords")),n.a.createElement("li",{key:"logout"},n.a.createElement(p.b,{to:"/logout"},"Logout"))],w=function(e){function t(e){var a;return Object(i.a)(this,t),(a=Object(u.a)(this,Object(m.a)(t).call(this,e))).sidebarRef=n.a.createRef(),a}return Object(d.a)(t,e),Object(l.a)(t,[{key:"componentDidMount",value:function(){var e=this;setTimeout(function(){window.M.Sidenav.init(e.sidebarRef.current)},0)}},{key:"componentDidUpdate",value:function(){window.M.Sidenav.getInstance(this.sidebarRef.current).close()}},{key:"render",value:function(){var e=this.props.services.authService.sessionExists(),t=e?f:v;return n.a.createElement("div",{className:"cp-navigation"},n.a.createElement("nav",null,n.a.createElement("div",{className:"nav-wrapper"},n.a.createElement(p.b,{to:"/",className:"brand-logo"},"Lorikeet"),n.a.createElement("a",{href:"#","data-target":"mobile-demo",className:"sidenav-trigger right hide-on-med-and-up"},n.a.createElement("i",{className:"material-icons"},"menu")),n.a.createElement("ul",{className:"right hide-on-small-and-down"},t))),n.a.createElement("ul",{className:"sidenav",id:"mobile-demo",ref:this.sidebarRef},t))}}]),t}(n.a.Component),E=Object(h.f)(w);a(55);function b(){return n.a.createElement("div",{className:"cp-home"},n.a.createElement("div",{className:"heading"},n.a.createElement("h1",null,"Lorikeet"),n.a.createElement("p",{className:"subtitle"},"A secure online password manager.")),n.a.createElement("div",{className:"bird-banner"},n.a.createElement("img",{alt:"",src:"".concat("","/bird_large.png")})),n.a.createElement("div",{className:"row"},n.a.createElement("div",{className:"col s12 m4"},n.a.createElement("div",{className:"center promo promo-example"},n.a.createElement("i",{className:"material-icons large"},"spa"),n.a.createElement("h5",{className:"promo-caption"},"Easy"),n.a.createElement("p",{className:"light center"},"You can stop keeping track of your passwords. It's easy to manage your passwords with Lorikeet."))),n.a.createElement("div",{className:"col s12 m4"},n.a.createElement("div",{className:"center promo promo-example"},n.a.createElement("i",{className:"material-icons large"},"vpn_key"),n.a.createElement("h5",{className:"promo-caption"},"Secure"),n.a.createElement("p",{className:"light center"},"With strong AES encryption on the client-side and server-side, you don't need to worry about your passwords."))),n.a.createElement("div",{className:"col s12 m4"},n.a.createElement("div",{className:"center promo promo-example"},n.a.createElement("i",{className:"material-icons large"},"favorite"),n.a.createElement("h5",{className:"promo-caption"},"Free"),n.a.createElement("p",{className:"light center"},"Lorikeet is free to use, and ",n.a.createElement("a",{href:"https://github.com/ClintonMorrison/lorikeet"},"open source"),". It was created with Golang and React. If you like it you can ",n.a.createElement("a",{href:"https://ko-fi.com/T6T0VOWY"},"support me on Ko-fi"),".")))),n.a.createElement(p.b,{to:"/register",className:"sign-up-link waves-effect waves-light btn-large btn"},"Sign Up Now"))}var y=a(17),g=a(7),k=a.n(g),N=function(e){function t(e){var a;return Object(i.a)(this,t),(a=Object(u.a)(this,Object(m.a)(t).call(this,e))).ref=n.a.createRef(),a}return Object(d.a)(t,e),Object(l.a)(t,[{key:"componentDidMount",value:function(){var e=this;setTimeout(function(){window.M.updateTextFields(),e.props.autoFocus&&e.ref.current.focus()},0)}},{key:"render",value:function(){var e=this.props,t=e.id,a=e.label,r=e.onChange,s=e.value,o=e.error,c=e.type,i=e.icon;return n.a.createElement("div",{className:"cp-text-field row"},n.a.createElement("div",{className:"input-field col s12"},i&&n.a.createElement("i",{className:"material-icons prefix"},i),n.a.createElement("input",{id:t,type:c,className:o?"invalid":"",autoComplete:this.props.autoComplete,value:s,onChange:function(e){return r(e.target.value)},ref:this.ref}),n.a.createElement("label",{htmlFor:t},a),n.a.createElement("span",{className:"helper-text","data-error":o})))}}]),t}(n.a.Component);N.defaultProps={type:"text"};var S=function(e){function t(e){var a;return Object(i.a)(this,t),(a=Object(u.a)(this,Object(m.a)(t).call(this,e))).state={username:"",password:"",usernameError:"",passwordError:""},a}return Object(d.a)(t,e),Object(l.a)(t,[{key:"submit",value:function(e){var t=this;e.preventDefault();var a=!0;if(this.state.username||(this.setState({usernameError:"Username cannot be empty"}),a=!1),this.state.password||(this.setState({passwordError:"Password cannot be empty"}),a=!1),a){var r=this.state,n=r.username,s=r.password;this.props.services.documentService.createDocument({username:n,password:s}).then(function(){t.props.history.push("/passwords")}).catch(function(e){console.log(Object(y.a)({},e));var a=k.a.get(e,"response.data.error","An error occurred.");a&&t.setState({usernameError:a})})}}},{key:"clearErrors",value:function(){this.setState({usernameError:"",passwordError:""})}},{key:"updateUsername",value:function(e){this.clearErrors(),this.setState({username:e})}},{key:"updatePassword",value:function(e){this.clearErrors(),this.setState({password:e})}},{key:"render",value:function(){var e=this;return n.a.createElement("div",{className:"cp-register"},n.a.createElement("h1",null,"Sign Up"),n.a.createElement("div",{className:"row"},n.a.createElement("form",{className:"col s12"},n.a.createElement("div",{className:"row"},n.a.createElement("div",{className:"col s12"},"Enter a username and a strong password for your new account.",n.a.createElement("p",null,n.a.createElement("strong",null,"Please write down your account information and keep it safe. ")),n.a.createElement("p",null,"Because of how your data will be encrypted, it will not be possible to regain control of your account if you forget."))),n.a.createElement(N,{label:"Username",id:"username",value:this.state.username,error:this.state.usernameError,onChange:function(t){return e.updateUsername(t)}}),n.a.createElement(N,{label:"Password",id:"password",type:"password",value:this.state.password,error:this.state.passwordError,onChange:function(t){return e.updatePassword(t)}}),n.a.createElement("div",{className:"row"},n.a.createElement("div",{className:"input-field col s12"},n.a.createElement("button",{className:"btn waves-effect waves-light",type:"submit",name:"action",onClick:function(t){return e.submit(t)}},"Register"))))))}}]),t}(n.a.Component),O=function(e){function t(e){var a;return Object(i.a)(this,t),(a=Object(u.a)(this,Object(m.a)(t).call(this,e))).state={username:"",password:"",usernameError:"",passwordError:""},a}return Object(d.a)(t,e),Object(l.a)(t,[{key:"submit",value:function(e){var t=this;e.preventDefault();var a=!0;if(this.state.username||(this.setState({usernameError:"Username cannot be empty"}),a=!1),this.state.password||(this.setState({passwordError:"Password cannot be empty"}),a=!1),a){var r=this.state,n=r.username,s=r.password;this.props.services.authService.setCredentials({username:n,password:s}),this.props.services.documentService.loadDocument().then(function(){t.props.history.push("/passwords")}).catch(function(e){console.log(e);var a=k.a.get(e,"response.data.error","An error occurred.");a&&t.setState({usernameError:a,passwordError:" "})})}}},{key:"clearErrors",value:function(){this.setState({usernameError:"",passwordError:""})}},{key:"updateUsername",value:function(e){this.clearErrors(),this.setState({username:e})}},{key:"updatePassword",value:function(e){this.clearErrors(),this.setState({password:e})}},{key:"render",value:function(){var e=this;return n.a.createElement("div",{className:"cp-login"},n.a.createElement("h1",null,"Login"),n.a.createElement("div",{className:"row"},n.a.createElement("form",{className:"col s12"},n.a.createElement("div",{className:"row"},n.a.createElement("div",{className:"col s12"},"Enter your username and password to login.")),n.a.createElement(N,{label:"Username",id:"username",autoComplete:"username",value:this.state.username,error:this.state.usernameError,onChange:function(t){return e.updateUsername(t)}}),n.a.createElement(N,{label:"Password",id:"password",autoComplete:"password",type:"password",value:this.state.password,error:this.state.passwordError,onChange:function(t){return e.updatePassword(t)}}),n.a.createElement("div",{className:"row"},n.a.createElement("div",{className:"input-field col s12"},n.a.createElement("button",{className:"btn waves-effect waves-light",type:"submit",name:"action",onClick:function(t){return e.submit(t)}},"Login"))))))}}]),t}(n.a.Component);function j(e){return e.services.authService.logout(),n.a.createElement(h.a,{to:"/login"})}a(57);var P=function(e){function t(e){var a;return Object(i.a)(this,t),(a=Object(u.a)(this,Object(m.a)(t).call(this,e))).state={oldPassword:"",newPassword:"",oldPasswordError:"",newPasswordError:""},a}return Object(d.a)(t,e),Object(l.a)(t,[{key:"submitPassword",value:function(e){var t=this;e.preventDefault(),this.clearErrors();var a=!1;this.state.oldPassword||(this.setState({oldPasswordError:"Old password cannot be empty"}),a=!0),this.state.newPassword||(this.setState({newPasswordError:"New password cannot be empty"}),a=!0),this.state.oldPassword||(this.setState({oldPasswordError:"Old password cannot be empty"}),a=!0),this.props.services.authService.passwordMatchesSession(this.state.oldPassword)||(this.setState({oldPasswordError:"Old password is incorrect"}),a=!0),a||this.props.services.documentService.updatePassword(this.state.newPassword).then(function(){t.props.history.push("/passwords")}).catch(function(e){var a=k.a.get(e,"response.data.error","An error occurred.");t.setState({newPasswordError:a})})}},{key:"submitDelete",value:function(e){var t=this;e.preventDefault(),window.confirm("Are you sure you want to delete your account? All your password data will be permanently deleted.")&&this.props.services.documentService.deleteDocument().then(function(){t.props.history.push("/logout")}).catch(function(e){var t=k.a.get(e,"response.data.error","An error occurred.");alert(t)})}},{key:"clearErrors",value:function(){this.setState({oldPasswordError:"",newPasswordError:""})}},{key:"updateOldPassword",value:function(e){this.clearErrors(),this.setState({oldPassword:e})}},{key:"updateNewPassword",value:function(e){this.clearErrors(),this.setState({newPassword:e})}},{key:"render",value:function(){var e=this;return n.a.createElement("div",{className:"cp-account"},n.a.createElement("h1",null,"Account"),n.a.createElement("div",{className:"row"},n.a.createElement("form",{className:"col s12"},n.a.createElement("h2",null,"Change Password"),n.a.createElement("p",null,n.a.createElement("strong",null,"Please write down your new password and keep it safe. ")),n.a.createElement("p",null,"Because of how your data will be encrypted, it will not be possible to regain control of your account if you forget."),n.a.createElement(N,{label:"Old Password",id:"old-passwrd",type:"password",value:this.state.oldPassword,error:this.state.oldPasswordError,onChange:function(t){return e.updateOldPassword(t)}}),n.a.createElement(N,{label:"New Password",id:"new-password",type:"password",value:this.state.newPassword,error:this.state.newPasswordError,onChange:function(t){return e.updateNewPassword(t)}}),n.a.createElement("div",{className:"row"},n.a.createElement("div",{className:"input-field col s12"},n.a.createElement("button",{className:"btn waves-effect waves-light",type:"submit",name:"action",onClick:function(t){return e.submitPassword(t)}},"Update Password"))))),n.a.createElement("div",{className:"row"},n.a.createElement("form",{className:"col s12"},n.a.createElement("h2",null,"Delete Account"),n.a.createElement("p",null,"If you no longer which you use Lorikeet to manage your passwords you can delete your account. This will delete all your stored passwords, and account data. You will not be able to login again, or view your passwords."),n.a.createElement("p",null,n.a.createElement("strong",null,"This is irreversible.")),n.a.createElement("div",{className:"row"},n.a.createElement("div",{className:"input-field col s12"},n.a.createElement("button",{className:"btn waves-effect waves-light btn-negative",type:"submit",name:"action",onClick:function(t){return e.submitDelete(t)}},"Delete All Data"))))))}}]),t}(n.a.Component),C=a(39);a(58);function D(e){var t=e.title,a=e.value,r=e.className,s=e.children;return a?n.a.createElement("div",{className:"cp-basic-field ".concat(r)},n.a.createElement("strong",null,t),": ",a,s):null}D.defaultProps={className:""};a(59);var T=function(e){function t(e){var a;return Object(i.a)(this,t),(a=Object(u.a)(this,Object(m.a)(t).call(this,e))).buttonRef=n.a.createRef(),a}return Object(d.a)(t,e),Object(l.a)(t,[{key:"copyToClipboard",value:function(e){e.preventDefault();var t=document.createElement("textarea");t.value=this.props.value,document.body.appendChild(t),t.select(),document.execCommand("copy"),document.body.removeChild(t),e.target.focus(),this.setState({copied:!0})}},{key:"getValue",value:function(){return this.props.mask?k.a.repeat("\u2022",this.props.value.length):this.props.value}},{key:"renderToggleLink",value:function(){var e=this;return n.a.createElement("button",{className:"copy-button btn-small waves-effect waves-light btn-negative",onClick:function(t){return e.copyToClipboard(t)},ref:this.buttonRef},n.a.createElement("i",{className:"material-icons left"},"content_copy"),this.props.title)}},{key:"render",value:function(){var e=this.getValue();return n.a.createElement("div",{className:"cp-copyable-field"},this.props.renderField&&n.a.createElement(D,{className:"cp-copyable-field",title:this.props.title,value:e}),this.renderToggleLink())}}]),t}(n.a.Component),x=(a(60),function(e){function t(){return Object(i.a)(this,t),Object(u.a)(this,Object(m.a)(t).apply(this,arguments))}return Object(d.a)(t,e),Object(l.a)(t,[{key:"renderTitle",value:function(){var e=this.props.item,t=e.title||"Untitled";return e.website?n.a.createElement("a",{target:"_blank",rel:"noopener noreferrer",href:e.website},t):t}},{key:"renderViewButton",value:function(){return n.a.createElement(p.b,{className:"waves-effect waves-light btn-small",to:"/passwords/".concat(this.props.item.id)},n.a.createElement("i",{className:"material-icons"},"edit"))}},{key:"render",value:function(){var e=this.props.item;return n.a.createElement("li",{className:"cp-item collection-item"},n.a.createElement("span",{className:"title"},this.renderTitle()),n.a.createElement(T,{title:"Username",value:e.username}),n.a.createElement(T,{title:"Password",value:e.password,mask:!0}),this.renderViewButton())}}]),t}(n.a.Component)),A=function(e){function t(){return Object(i.a)(this,t),Object(u.a)(this,Object(m.a)(t).apply(this,arguments))}return Object(d.a)(t,e),Object(l.a)(t,[{key:"render",value:function(){return n.a.createElement("ul",{className:"collection"},this.props.passwords.map(function(e){return n.a.createElement(x,{key:e.id||e.title,item:e})}))}}]),t}(n.a.Component);a(61);function U(){return n.a.createElement("div",{className:"cp-loader"},n.a.createElement("div",{className:"preloader-wrapper big active"},n.a.createElement("div",{className:"spinner-layer spinner-blue-only"},n.a.createElement("div",{className:"circle-clipper left"},n.a.createElement("div",{className:"circle"})),n.a.createElement("div",{className:"gap-patch"},n.a.createElement("div",{className:"circle"})),n.a.createElement("div",{className:"circle-clipper right"},n.a.createElement("div",{className:"circle"})))))}a(62);var I=a(37),M=a.n(I),L=function(e){function t(e){var a;return Object(i.a)(this,t),(a=Object(u.a)(this,Object(m.a)(t).call(this,e))).state={query:"",document:null},a}return Object(d.a)(t,e),Object(l.a)(t,[{key:"applyFilter",value:function(e){var t=this.state.query.toLowerCase().trim();if(!t)return!0;return["title","notes","website"].some(function(a){return-1!==e[a].toLowerCase().trim().indexOf(t)})}},{key:"getPasswords",value:function(){var e=this;return k.a.chain(this.state).get("document.passwords",[]).filter(function(t){return e.applyFilter(t)}).value()}},{key:"createPassword",value:function(e){var t=this;e.preventDefault();var a=M()(),r={id:a,title:"",username:"",password:"",email:"",website:"",notes:""};this.props.services.documentService.loadDocument().then(function(e){e.passwords=[].concat(Object(C.a)(e.passwords),[r]),t.props.services.documentService.updateDocument(e)}).then(function(){t.props.history.push("/passwords/".concat(a))}).catch(function(){t.props.history.push("/logout")})}},{key:"componentDidMount",value:function(){var e=this;return this.props.services.documentService.loadDocument().then(function(t){e.setState({document:t})}).catch(function(){e.props.history.push("/logout")})}},{key:"render",value:function(){var e=this;return this.state.document?n.a.createElement("div",{className:"cp-passwords"},n.a.createElement("h1",null,"Passwords"),n.a.createElement("form",{onSubmit:function(e){return e.preventDefault()}},n.a.createElement(N,{label:"Search",id:"search",icon:"search",value:this.state.query,onChange:function(t){return e.setState({query:t})}})),n.a.createElement("div",{className:"top-actions"},n.a.createElement("button",{onClick:function(t){return e.createPassword(t)},className:"waves-effect waves-light btn"},n.a.createElement("i",{className:"material-icons left"},"add"),"Add")),n.a.createElement(A,{passwords:this.getPasswords()})):n.a.createElement(U,null)}}]),t}(n.a.Component),R=a(18),F=function(e){function t(){return Object(i.a)(this,t),Object(u.a)(this,Object(m.a)(t).apply(this,arguments))}return Object(d.a)(t,e),Object(l.a)(t,[{key:"componentDidMount",value:function(){setTimeout(function(){window.M.updateTextFields(),window.M.textareaAutoResize(window.$("textarea"))},0)}},{key:"render",value:function(){var e=this.props,t=e.id,a=e.label,r=e.onChange,s=e.value,o=e.error,c=e.icon;return n.a.createElement("div",{className:"cp-text-field row"},n.a.createElement("div",{className:"input-field col s12"},c&&n.a.createElement("i",{className:"material-icons prefix"},c),n.a.createElement("textarea",{id:t,className:"materialize-textarea ".concat(o?"invalid":""),value:s,onChange:function(e){return r(e.target.value)}}),n.a.createElement("label",{htmlFor:t},a),n.a.createElement("span",{className:"helper-text","data-error":o})))}}]),t}(n.a.Component),W=(a(65),function(e){function t(e){var a;return Object(i.a)(this,t),(a=Object(u.a)(this,Object(m.a)(t).call(this,e))).state={item:{title:"",username:"",password:"",email:"",website:"",notes:""},errors:{}},a}return Object(d.a)(t,e),Object(l.a)(t,[{key:"updateItem",value:function(e,t){this.setState({item:Object(y.a)({},this.state.item,Object(R.a)({},e,t))})}},{key:"componentDidMount",value:function(){this.setState({item:Object(y.a)({},this.props.item)})}},{key:"handleSave",value:function(e){e.preventDefault(),this.props.updatePassword(this.state.item)}},{key:"handleDelete",value:function(e){if(e.preventDefault(),window.confirm("Are you sure you want to delete this password? This cannot be undone."))return this.props.deletePassword(this.state.item)}},{key:"render",value:function(){var e=this,t=this.state.item;return t?n.a.createElement("div",{className:"cp-details"},n.a.createElement("h1",{className:"title"},t.title||"Untitled"),n.a.createElement("form",null,n.a.createElement(N,{autoFocus:!0,label:"Title",id:"title",icon:"title",value:t.title,error:this.state.errors.title,onChange:function(t){return e.updateItem("title",t)}}),n.a.createElement(N,{label:"Username",id:"username",icon:"person",value:t.username,error:this.state.errors.username,onChange:function(t){return e.updateItem("username",t)}}),n.a.createElement(N,{label:"Password",id:"password",type:"text",icon:"vpn_key",value:t.password,error:this.state.errors.password,onChange:function(t){return e.updateItem("password",t)}}),n.a.createElement(N,{label:"Email",id:"email",type:"text",icon:"email",value:t.email,error:this.state.errors.email,onChange:function(t){return e.updateItem("email",t)}}),n.a.createElement(N,{label:"Website",id:"website",type:"text",icon:"cloud",value:t.website,error:this.state.errors.website,onChange:function(t){return e.updateItem("website",t)}}),n.a.createElement(F,{label:"Notes",id:"notes",type:"text",icon:"note",value:t.notes,error:this.state.errors.notes,onChange:function(t){return e.updateItem("notes",t)}}),n.a.createElement("div",{className:"actions"},n.a.createElement(p.b,{className:"waves-effect waves-light grey btn",to:"/passwords"},"Cancel"),n.a.createElement("button",{className:"waves-effect waves-light btn",onClick:function(t){return e.handleSave(t)}},"Save"),n.a.createElement("button",{className:"waves-effect waves-light btn btn-negative delete-button",onClick:function(t){return e.handleDelete(t)}},n.a.createElement("i",{className:"material-icons"},"delete"))))):null}}]),t}(n.a.Component)),B=function(e){function t(e){var a;return Object(i.a)(this,t),(a=Object(u.a)(this,Object(m.a)(t).call(this,e))).state={document:null},a}return Object(d.a)(t,e),Object(l.a)(t,[{key:"updatePassword",value:function(e){var t=this,a=e.id;this.props.services.documentService.loadDocument().then(function(r){var n=k.a.findIndex(r.passwords,{id:a});r.passwords[n]=e,t.props.services.documentService.updateDocument(r)}).then(function(){t.props.history.push("/passwords")}).catch(function(){t.props.history.push("/logout")})}},{key:"deletePassword",value:function(e){var t=this,a=e.id;this.props.services.documentService.loadDocument().then(function(e){var r=k.a.findIndex(e.passwords,{id:a});e.passwords.splice(r,1),t.props.services.documentService.updateDocument(e)}).then(function(){t.props.history.push("/passwords")}).catch(function(){t.props.history.push("/logout")})}},{key:"getPasswords",value:function(){return k.a.get(this.state,"document.passwords",[])}},{key:"componentDidMount",value:function(){var e=this;return this.props.services.documentService.loadDocument().then(function(t){e.setState({document:t})}).catch(function(){e.props.history.push("/logout")})}},{key:"render",value:function(){var e=this;if(!this.state.document)return n.a.createElement(U,null);var t=this.props.match.params.id,a=this.getPasswords(),r=k.a.find(a,{id:t});return r?n.a.createElement("div",{className:"cp-view"},n.a.createElement(W,{item:r,createPassword:function(){return e.createPassword()},updatePassword:function(t){return e.updatePassword(t)},deletePassword:function(t){return e.deletePassword(t)}})):n.a.createElement(h.a,{to:"/passwords"})}}]),t}(n.a.Component),H=function(e){function t(){return Object(i.a)(this,t),Object(u.a)(this,Object(m.a)(t).apply(this,arguments))}return Object(d.a)(t,e),Object(l.a)(t,[{key:"renderPage",value:function(e,t){return n.a.createElement(e,Object.assign({services:this.props.services},t))}},{key:"render",value:function(){var e=this,t=this.props.services;return n.a.createElement(p.a,null,n.a.createElement("header",null,n.a.createElement(E,{services:t})),n.a.createElement("main",{className:"container"},n.a.createElement(h.d,null,n.a.createElement(h.b,{path:"/",exact:!0,render:function(t){return e.renderPage(b,t)}}),n.a.createElement(h.b,{path:"/login",exact:!0,render:function(t){return e.renderPage(O,t)}}),n.a.createElement(h.b,{path:"/logout",exact:!0,render:function(t){return e.renderPage(j,t)}}),n.a.createElement(h.b,{path:"/register",exact:!0,render:function(t){return e.renderPage(S,t)}}),n.a.createElement(h.b,{path:"/account",exact:!0,render:function(t){return e.renderPage(P,t)}}),n.a.createElement(h.b,{path:"/passwords",exact:!0,render:function(t){return e.renderPage(L,t)}}),n.a.createElement(h.b,{path:"/passwords/:id",render:function(t){return e.renderPage(B,t)}}),n.a.createElement(h.b,{render:function(){return n.a.createElement(h.a,{to:"/"})}}))))}}]),t}(n.a.Component),J=a(23),V=a.n(J),_=a(24),q=a.n(_),z=a(38),Y=a.n(z),$=function(){function e(t){var a=t.apiService;Object(i.a)(this,e),this.apiService=a,this.document=null}return Object(l.a)(e,[{key:"hashPassword",value:function(e){return V()(e+"CC352C99A14616AD22678563ECDA5").toString()}},{key:"hashToken",value:function(e){return V()(e+"7767B9225CF66B418DD2A39CBC4AA").toString()}},{key:"passwordMatchesSession",value:function(e){return e&&this.hashPassword(e)===this.getToken()}},{key:"setCredentials",value:function(e){var t=e.username,a=e.password;sessionStorage.setItem("username",t),this.setToken(a)}},{key:"setToken",value:function(e){sessionStorage.setItem("token",this.hashPassword(e))}},{key:"sessionExists",value:function(){return!(!this.getUsername()||!this.getToken())}},{key:"getUsername",value:function(){return sessionStorage.getItem("username")}},{key:"getToken",value:function(){return sessionStorage.getItem("token")}},{key:"logout",value:function(){sessionStorage.clear()}},{key:"getHashedToken",value:function(){var e=this.getToken();return e?this.hashToken(e):null}},{key:"encryptWithToken",value:function(e){var t=arguments.length>1&&void 0!==arguments[1]&&arguments[1];return t=t||this.getToken(),q.a.encrypt(e,t).toString()}},{key:"decryptWithToken",value:function(e){var t=this.getToken();return q.a.decrypt(e,t).toString(Y.a)}},{key:"getHeaders",value:function(){var e=this.getUsername(),t=this.getHashedToken(),a=btoa("".concat(e,":").concat(t));return{Authorization:"Basic ".concat(a)}}}]),e}(),G=a(15),K=a.n(G),Q=function(){function e(t){var a=t.baseURL,r=t.authService;Object(i.a)(this,e),this.authService=r,K.a.defaults.baseURL=a,K.a.defaults.headers.common.Accept="application/json"}return Object(l.a)(e,[{key:"get",value:function(e,t,a){return K()({method:"get",url:"/".concat(e),params:t,headers:a})}},{key:"post",value:function(e,t,a){return K()({method:"post",url:"/".concat(e),data:t,headers:a})}},{key:"put",value:function(e,t,a){return K()({method:"put",url:"/".concat(e),data:t,headers:a})}},{key:"del",value:function(e,t){return K()({method:"delete",url:"/".concat(e),headers:t})}}]),e}(),X=function(){function e(t){var a=t.apiService,r=t.authService;Object(i.a)(this,e),this.apiService=a,this.authService=r}return Object(l.a)(e,[{key:"createDocument",value:function(e){var t=e.username,a=e.password;this.authService.setCredentials({username:t,password:a});var r=JSON.stringify({passwords:[]}),n=this.authService.encryptWithToken(r);return this.apiService.post("document",n,this.authService.getHeaders())}},{key:"loadDocument",value:function(){var e=this;return this.apiService.get("document",{},this.authService.getHeaders()).then(function(t){var a=k.a.get(t,"data.document")||"{}",r=e.authService.decryptWithToken(a);return e.document=JSON.parse(r),e.document})}},{key:"updateDocument",value:function(e){var t=JSON.stringify(e),a=this.authService.encryptWithToken(t);return this.apiService.put("document",a,this.authService.getHeaders())}},{key:"deleteDocument",value:function(){return this.apiService.del("document",this.authService.getHeaders())}},{key:"updatePassword",value:function(e){var t=this,a=this.authService.hashPassword(e),r=this.authService.hashToken(a);return this.loadDocument().then(function(n){var s=JSON.stringify(n),o=t.authService.encryptWithToken(s,a);return t.apiService.put("document/password",{password:r,document:o},t.authService.getHeaders()).then(function(){return t.authService.setToken(e)})})}}]),e}(),Z=(a(90),new Q({baseURL:"".concat(window.location.origin,"/api/")})),ee=new $({apiService:Z}),te={apiService:Z,authService:ee,documentService:new X({authService:ee,apiService:Z})};window.services=te;var ae=function(){return n.a.createElement("div",{className:"cp-app"},n.a.createElement(H,{services:te}),n.a.createElement(c,null))};Boolean("localhost"===window.location.hostname||"[::1]"===window.location.hostname||window.location.hostname.match(/^127(?:\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$/));o.a.render(n.a.createElement(ae,null),document.getElementById("root")),"serviceWorker"in navigator&&navigator.serviceWorker.ready.then(function(e){e.unregister()})}},[[40,1,2]]]);
//# sourceMappingURL=main.5fd2d6cc.chunk.js.map