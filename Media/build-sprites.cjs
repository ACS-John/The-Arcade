// Original editable vector game sprites; PNGs are the native BR runtime build.
// Run with Node and sharp available (NODE_PATH may point to the shared runtime).
const fs=require('fs'),path=require('path'),sharp=require('sharp');
const out=path.join(__dirname,'../art/sprites');fs.mkdirSync(out,{recursive:true});
const colors=['#50f2ec','#ff63ce','#ffe37b','#81fa9c','#ff687c','#ffffff','#658fff'];
const tile=(body,c='#50f2ec',bg='#080c21')=>`<svg xmlns="http://www.w3.org/2000/svg" width="96" height="96" viewBox="0 0 96 96"><defs><linearGradient id="metal" x2=".35" y2="1"><stop stop-color="#ffffff"/><stop offset=".28" stop-color="${c}"/><stop offset="1" stop-color="#17253d"/></linearGradient><radialGradient id="orb"><stop stop-color="#fff"/><stop offset=".25" stop-color="${c}"/><stop offset="1" stop-color="${c}" stop-opacity="0"/></radialGradient></defs><rect width="96" height="96" fill="${bg}"/>${body}</svg>`;
const files={};const add=(n,b,c,bg)=>files[n]=tile(b,c,bg);
add('empty','');
add('floor','<path d="M0 95H96M95 0V96" stroke="#13203a" opacity=".38"/>');
add('wall','<rect x="3" y="3" width="90" height="90" rx="15" fill="#0d1834" stroke="#1b3b73" stroke-width="5"/><rect x="9" y="9" width="78" height="78" rx="11" fill="#111e3c" stroke="#418ddd" stroke-width="2"/><path d="M17 71V24Q17 17 24 17H71" fill="none" stroke="#74caff" stroke-width="2" opacity=".65"/>');
add('pellet','<circle cx="48" cy="48" r="20" fill="url(#orb)"/><circle cx="48" cy="48" r="7" fill="#fff1b5"/>','#ffe37b');
add('power','<circle cx="48" cy="48" r="42" fill="url(#orb)"/><circle cx="48" cy="48" r="21" fill="none" stroke="#ffe37b" stroke-width="4"/><path d="M48 28L54 41L68 48L54 55L48 68L42 55L28 48L42 41Z" fill="#fff6d7"/>','#ffe37b');
add('ball','<circle cx="48" cy="48" r="43" fill="url(#orb)"/><circle cx="48" cy="48" r="17" fill="url(#metal)"/><circle cx="43" cy="42" r="5" fill="white"/>');
add('laser','<ellipse cx="48" cy="48" rx="21" ry="48" fill="url(#orb)"/><rect x="44" y="4" width="8" height="85" rx="4" fill="#fff5c9"/>','#ffe37b');
add('bomb','<path d="M31 13L61 33L40 52L66 73L44 91" fill="none" stroke="#ff576e" stroke-width="12"/><path d="M31 13L61 33L40 52L66 73L44 91" fill="none" stroke="#ffe1dc" stroke-width="3"/>','#ff687c');
add('ship','<ellipse cx="48" cy="77" rx="19" ry="18" fill="url(#orb)"/><path d="M48 8L58 43L85 72L60 67L48 57L36 67L11 72L38 43Z" fill="url(#metal)" stroke="#a2faff" stroke-width="2"/><path d="M48 24L55 51H41Z" fill="#163952" stroke="#fff"/><path d="M38 68L43 88M58 68L53 88" stroke="#ff9459" stroke-width="5"/>');
for(let i=0;i<7;i++){
 let c=colors[i];
 add('alien'+i,'<path d="M26 28L17 11M70 28L79 11M21 57L9 72M75 57L87 72" stroke="'+c+'" stroke-width="6"/><path d="M17 43L30 22H66L79 43L71 72H25Z" fill="url(#metal)" stroke="'+c+'" stroke-width="3"/><path d="M29 41L43 46L36 56L25 50M67 41L53 46L60 56L71 50" fill="#060b19"/><path d="M37 66H59" stroke="#fff" stroke-width="3"/>',c);
 add('drone'+i,'<circle cx="48" cy="46" r="34" fill="url(#metal)" stroke="'+c+'" stroke-width="3"/><path d="M20 68L15 86L35 76L48 87L61 76L81 86L76 68" fill="'+c+'"/><rect x="23" y="30" width="50" height="27" rx="13" fill="#080c21"/><ellipse cx="37" cy="42" rx="5" ry="8" fill="#fff"/><ellipse cx="59" cy="42" rx="5" ry="8" fill="#fff"/>',c);
 add('brick'+i,'<rect x="3" y="12" width="90" height="72" rx="10" fill="url(#metal)" stroke="'+c+'" stroke-width="3"/><path d="M12 69V23H83" fill="none" stroke="#fff" stroke-opacity=".65" stroke-width="3"/><path d="M18 77H79L88 69" fill="none" stroke="#000" stroke-opacity=".5" stroke-width="5"/>',c);
 add('block'+(i+1),'<rect x="3" y="3" width="90" height="90" rx="11" fill="url(#metal)" stroke="'+c+'" stroke-width="3"/><path d="M12 76V13H78" fill="none" stroke="#fff" stroke-opacity=".7" stroke-width="3"/><rect x="23" y="23" width="50" height="50" rx="7" fill="'+c+'" opacity=".28"/><path d="M19 83H80L86 77" fill="none" stroke="#000" stroke-opacity=".5" stroke-width="5"/>',c);
 add('car'+i,'<rect x="21" y="5" width="54" height="86" rx="15" fill="url(#metal)" stroke="'+c+'" stroke-width="3"/><path d="M26 25L31 17H65L70 25V38H26ZM28 66H68V79H28Z" fill="#0d1c35"/><rect x="15" y="23" width="7" height="16" rx="2" fill="#182439"/><rect x="74" y="23" width="7" height="16" rx="2" fill="#182439"/><path d="M26 10H36M60 10H70" stroke="#fff4bf" stroke-width="5"/><path d="M26 86H35M61 86H70" stroke="#ff4f58" stroke-width="4"/>',c);
}
add('ghost','<rect x="5" y="5" width="86" height="86" rx="10" fill="#111e36" stroke="#629aaf" stroke-dasharray="9 5" stroke-width="3"/>');
for(let d=1;d<=4;d++)for(let frame=0;frame<2;frame++){
 const angle={1:180,2:0,3:270,4:90}[d];
 add('munch'+d+'_'+frame,`<g transform="rotate(${angle} 48 48)"><path d="${frame?'M48 48L84 22A43 43 0 1 0 84 74Z':'M48 48L90 42A43 43 0 1 0 90 54Z'}" fill="url(#metal)" stroke="#ffe7a1" stroke-width="2"/><circle cx="50" cy="23" r="5" fill="#111a33"/></g>`,'#ffe37b');
}
add('shield','<path d="M10 83V39Q10 11 48 11Q86 11 86 39V83H66V66Q48 47 30 66V83Z" fill="url(#metal)" stroke="#91fcb3" stroke-width="3"/>','#81fa9c');
add('saucer','<ellipse cx="48" cy="55" rx="43" ry="18" fill="url(#metal)"/><path d="M24 48Q28 8 48 8Q68 8 72 48Z" fill="#382451" stroke="#ff92dc" stroke-width="3"/><path d="M13 59H83" stroke="#fff" stroke-width="3"/><circle cx="26" cy="62" r="4" fill="#ff63ce"/><circle cx="48" cy="66" r="4" fill="#ff63ce"/><circle cx="70" cy="62" r="4" fill="#ff63ce"/>','#ff63ce');
add('paddle','<rect x="0" y="24" width="96" height="50" rx="8" fill="url(#metal)" stroke="#68fff1" stroke-width="3"/><path d="M4 34H92" stroke="#fff" stroke-width="4"/><path d="M4 63H92" stroke="#144454" stroke-width="7"/>');
add('wide','<rect x="16" y="10" width="64" height="76" rx="24" fill="url(#metal)"/><path d="M31 36L21 48L31 60M65 36L75 48L65 60M22 48H74" fill="none" stroke="#081c2a" stroke-width="6"/>','#81fa9c');
add('multi','<rect x="16" y="10" width="64" height="76" rx="24" fill="url(#metal)"/><circle cx="48" cy="34" r="10" fill="#fff"/><circle cx="34" cy="58" r="10" fill="#fff"/><circle cx="62" cy="58" r="10" fill="#fff"/>','#ff63ce');
add('water','<path d="M-9 25Q12 8 36 25T84 25T132 25M-25 62Q0 45 24 62T72 62T120 62" fill="none" stroke="#247a9b" stroke-width="4"/><path d="M0 34H23M53 75H86" stroke="#58c7ce" stroke-width="2" opacity=".5"/>','#50f2ec','#102d43');
add('road','<path d="M0 47H26M70 47H96" stroke="#b7c4ce" opacity=".35" stroke-width="3"/>','#fff','#192537');
add('log','<rect y="6" width="96" height="84" rx="5" fill="#754e37"/><path d="M0 15H96M0 81H96" stroke="#c9945c" stroke-width="5"/><path d="M0 31Q22 22 43 34T96 31M0 62Q35 72 58 58T96 62M3 46H88" stroke="#3c302e" stroke-width="3" fill="none"/><ellipse cx="39" cy="42" rx="16" ry="9" fill="none" stroke="#ae784d" stroke-width="3"/>');
add('grass','<path d="M4 94L13 32L21 78L32 13L45 72L53 31L69 88L85 24L92 95" fill="#295d4f"/><path d="M0 95L26 54L30 93L61 59L59 96L85 72L96 96" fill="#429575"/>','#81fa9c','#123732');
add('home','<ellipse cx="48" cy="52" rx="40" ry="31" fill="#21614c" stroke="#74d39a" stroke-width="3"/><path d="M48 22L48 55L72 28" fill="#103d36"/>','#81fa9c','#102d43');
add('frog','<path d="M30 24L13 13L7 26L23 42M66 24L83 13L89 26L73 42M28 65L12 82L24 89L39 72M68 65L84 82L72 89L57 72" stroke="#61de8b" stroke-width="10" fill="none"/><ellipse cx="48" cy="49" rx="28" ry="32" fill="url(#metal)"/><circle cx="32" cy="23" r="13" fill="#93fcae"/><circle cx="64" cy="23" r="13" fill="#93fcae"/><circle cx="32" cy="21" r="6" fill="#122135"/><circle cx="64" cy="21" r="6" fill="#122135"/><path d="M33 52Q48 63 63 52" fill="none" stroke="#224836" stroke-width="3"/>','#81fa9c');
add('mushroom','<path d="M38 48L33 87H65L58 48" fill="#f2cdeb"/><path d="M9 51Q12 5 48 6Q84 5 87 51Q49 67 9 51Z" fill="url(#metal)" stroke="#ff9add" stroke-width="3"/><ellipse cx="33" cy="31" rx="8" ry="6" fill="#fff4fe"/><ellipse cx="62" cy="36" rx="9" ry="7" fill="#fff4fe"/><path d="M19 52Q48 64 78 52" stroke="#c1418a" stroke-width="3" fill="none"/>','#ff63ce');
add('worm','<path d="M25 24L8 12M20 47L4 47M25 72L8 87M71 24L88 12M76 47L92 47M71 72L88 87" stroke="#8befa7" stroke-width="5"/><ellipse cx="48" cy="48" rx="30" ry="39" fill="url(#metal)" stroke="#bcffd1" stroke-width="3"/><path d="M21 38Q48 52 75 38M21 57Q48 71 75 57" fill="none" stroke="#1d6c4f" stroke-width="4"/>','#81fa9c');
add('spider','<path d="M35 34L14 11L6 33M32 44L6 40L1 59M35 57L14 74L6 91M61 34L82 11L90 33M64 44L90 40L95 59M61 57L82 74L90 91" fill="none" stroke="#ef718d" stroke-width="5"/><ellipse cx="48" cy="48" rx="24" ry="30" fill="url(#metal)"/><circle cx="40" cy="34" r="4" fill="#fff"/><circle cx="56" cy="34" r="4" fill="#fff"/>','#ff687c');
// Larger ships/formation enemies span three logical cells, matching their hit area.
// Slice the vector viewBox at build time; this never modifies generated key art.
for(const name of ['ship',...Array.from({length:7},(_,i)=>'alien'+i),...Array.from({length:7},(_,i)=>'brick'+i)]){
 for(let part=0;part<3;part++)files[name+'-'+part]=files[name].replace('viewBox="0 0 96 96"',`preserveAspectRatio="none" viewBox="${part*32} 0 32 96"`);
}
(async()=>{for(const [name,svg] of Object.entries(files)){fs.writeFileSync(path.join(out,name+'.svg'),svg);await sharp(Buffer.from(svg)).png().toFile(path.join(out,name+'.png'));}console.log('Built '+Object.keys(files).length+' editable vector sprites and runtime PNGs.');})().catch(e=>{console.error(e);process.exit(1)});
