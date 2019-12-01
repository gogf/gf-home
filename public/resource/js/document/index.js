var currentUri = window.location.pathname;

// 返回/前进浏览器事件
window.onpopstate = function() {
    if (currentUri != window.location.pathname) {
        loadMarkdown(window.location.pathname, false)
    }
};

// 打开节点,参数为li节点的jquery对象
function openNode(node) {
    if (node.find("ul").eq(0).css('display') == "none") {
        node.find("ul").eq(0).show();
        node.find("i").eq(0).attr("class", "am-icon-caret-down");
    }
}

// 替换a标签及img标签的地址，加上/前缀修改为绝对路径
function replaceHrefAndSrc() {
    // 修改a/img标签链接，给相对路径统一加上前缀
    $(document).find("a").each(function(){
        var href = $(this).attr("href");
        if (typeof href != "undefined" && href.length > 0) {
            if (href.substr(0, 7) == "mailto:") {
                return
            }
            if (href.substr(0, 1) != "#" && href.substr(0, 4) != "http" && href.substr(0, 10) != "javascript") {
                if (href.substr(0, 1) != "/") {
                    href = "/" + href;
                }
                if (href.indexOf(".md") != -1) {
                    href = href.replace(".md", "");
                    $(this).attr("href", "javascript:loadMarkdown('" + href + "', true);");
                } else {
                    $(this).attr("href", href);
                }
            }
            if (href.substr(0, 4) == "http") {
                $(this).attr("target", "_blank");
            }
        }
    });
    $(document).find("img").each(function(){
        var src = $(this).attr("src");
        if (typeof src != "undefined" && src.length > 0) {
            if (src.substr(0, 1) != "/" && src.substr(0, 4) != "http") {
                $(this).attr("src", "/" + src);
            }
        }
    });
}

// 取消当前的li高亮
function cancelAllHighlight() {
    $("#side-markdown-view").find("li").removeClass("active");
}

// 高亮并展开当前打开的地址
function highlightLiByUri(uri) {
    var seli = $("a[href='"+ uri +"']").parent("li");
    seli.addClass("active");
    // 层级打开节点
    $("a[href='"+ uri +"']").parents("li").each(function(){
        openNode($(this));
    });

    var menu = $("#side-markdown-view");
    var windowHeiget = window.innerHeight;
    var menuScrollTop = menu.scrollTop();
    var seliTop = seli.offset().top;
    // console.log(windowHeiget, seliTop, menuScrollTop);
    // 将菜单滚动调整到当前页面的菜单所在位置
    if (seliTop + 46 > windowHeiget) {
        // 菜单所在位置在页面以下
        menu.animate({ scrollTop: menuScrollTop + seliTop - 8 }, "slow");
    } else if (seliTop < 0) {
        // 菜单所在位置在页面以上
        menu.animate({ scrollTop: menuScrollTop + seliTop - 8 }, "slow");
    }
}

// 监听按钮事件监听
function copyBtnOn() {
    $('.copy-code').on('click',function() {
        var span=$(this);
        var id=span.attr("code-id");
        var codeContent=$("#code-content-id-"+id);
        if(copyText(codeContent.text())){
            //span.css("color","#00ff00");
            span.html(`<i class="doc-act-clip am-icon-copy"></i>done`);
        }else{
            //span.css("color","red");
            span.html(`<i class="doc-act-clip am-icon-copy"></i>fail`);
        }
        setTimeout(function(){
            //span.css("color","");
            span.html(`<i class="doc-act-clip am-icon-copy"></i>copy`);
        },500);
    });
}
// 复制功能
function copyText(text) {
    var textarea = document.createElement("textarea");//创建input对象
    var currentFocus = document.activeElement;//当前获得焦点的元素
    document.body.appendChild(textarea);//添加元素
    textarea.value = text;
    textarea.focus();
    if(textarea.setSelectionRange)
        textarea.setSelectionRange(0, textarea.value.length);//获取光标起始位置到结束位置
    else
        textarea.select();
    try {
        var flag = document.execCommand("copy");//执行复制
    } catch(eo) {
        var flag = false;
    }
    document.body.removeChild(textarea);//删除元素
    currentFocus.focus();
    return flag;
}
// TOC事件
function tocOn() {
    if (navigator.userAgent.lastIndexOf("Mobile") == -1) {
        $("div.toc").addClass("toc-pc");
        $("div.toc").hover(function () {
            $(this).css({ opacity: 1 });
        }, function () {
            $(this).css({ opacity: .2 });
        });
    }
}
// 插入代码
function isEleExist(id) {
    if($("#"+id).length <= 0) {
        $("body").append($("<div>").attr("id",id).hide());     
    }
}
// 判断元素是否存在滚动条
function hasScrolled(element,direction){
    if(direction==='vertical'){
        return element.scrollHeight>element.clientHeight;
    }else if(direction==='horizontal'){
        return element.scrollWidth>element.clientWidth;
    }
}

// 重新解析markdown内容
function reloadMainMarkdown() {
    var content = $("#main-markdown-content").text()
    if (content.length > 0) {
        isEleExist("code-list");
        $("#code-list").html("");
        $('#main-markdown-view').html(marked($("#main-markdown-content").text()));
        $('#main-markdown-view pre code').each(function(i, block) {
            var thisBlock=$(block);
            //记录代码块内容
            var codeContent=$("<span>").text(thisBlock.text()).attr("id","code-content-id-"+i);
            $("#code-list").append(codeContent);
            // 添加复制按钮，添加class用于事件监听
            var copyBtn=$("<span>").attr({
                "style":"position:absolute;right:0px;top:0px;cursor:pointer;user-select:none;padding: 2px 8px;font-size:14px;",
                "title":"copy",
                "code-id":""+i
            }).addClass("copy-code");
            copyBtn.html(`<i class="doc-act-clip am-icon-copy"></i>copy`);
            var copyDiv = $("<div>").attr({
                "style":"color:#f8f8f2;position:relative; z-index:999;margin-top: 8px;"
            }); 
            copyDiv.append(copyBtn)
            thisBlock.parent().before(copyDiv);
            thisBlock.parent().attr("style","position: relative;").attr("class","check-scroll");
            Prism.highlightElement(block);
            
        });
        //用于检测代码块是否有纵向滚动条
        $(".check-scroll").each(function(){
            var ua=navigator.userAgent;
            if(hasScrolled(this ,'vertical') && ua.search("Mobile")==-1){
                if(ua.search("Windows NT")!=-1){
                    $(this).prev().find("span").css("padding","2px 22px");
                }else{
                    $(this).prev().find("span").css("padding","2px 16px");
                }
            }
        });
        // 生成TOC菜单
        $('#main-markdown-toc').html("");
        new Toc('main-markdown-view', {
            'level'   : 4,
            'class'   : 'toc',
            'targetId': 'main-markdown-toc'
        } );
        if ($('#main-markdown-toc').html().length > 0) {
            var html = $("#main-markdown-view").html().replace("<p>[TOC]</p>", $('#main-markdown-toc').html());
            html += $("#powered").html();
            $("#main-markdown-view").html(html)
        }

        copyBtnOn();
        tocOn();
    }
    replaceHrefAndSrc();
    updateHelpUrl(window.location.pathname);
}

// 更新文档markdown链接地址
function updateHelpUrl(uri) {
    $("#help-icon").attr("href", "https://github.com/gogf/gf-doc/tree/master" + uri + ".md");
}

// 修改当前标题
function updateWindowTitle(uri) {
    var title = "";
    $("a[href='"+ uri +"']").parents("li").each(function(){
        if (title == "") {
            title  = $(this).find("a").eq(0).text()
        } else {
            title += " - " + $(this).find("a").eq(0).text()
        }
    });
    document.title = title + " - " + baseTitle;
}

// 请求markdown内容
function loadMarkdown(uri, addState) {
    currentUri = uri;
    cancelAllHighlight();
    // 添加历史记录
    if (addState) {
        window.history.pushState({
            title : document.title,
            uri   : window.location.pathname
        }, document.title, window.location.origin + uri);
    }
    highlightLiByUri(uri);
    updateWindowTitle(uri);
    $("#main-markdown-view").html("<div class=\"loading-small\"></div> Loading...");
    // AJAX读取文档
    $.ajax({
        type     : "get",
        url      : uri,
        dataType : "json",
        success: function(result){
            if (result.code == 1) {
                $("#main-markdown-content").text(result.data);
                reloadMainMarkdown();
            }
        }
    });
}

$(function() {
    reloadMainMarkdown();
    // 修改list样式
    $("#side-markdown-view").find("ul").addClass("am-list am-list-border");
    // 回到顶部
    $('#totop-icon').on('click', function() {
        $("#main-markdown-view").smoothScroll({position: 0, speed: 300});
    });
    $("#main-markdown-view").scroll(function() {
        if ($(this).scrollTop() >= 400) {
            $("#totop-icon").show();
        } else {
            $("#totop-icon").hide();
        }
    });
    // 菜单树形结构处理
    function indentTreeByLi(obj) {
        if ($(obj).find("ul").length > 0) {
            // 添加图标
            $(obj).find("a").eq(0).find("i").attr("class", "am-icon-caret-down");
            // 递归处理树型缩进
            $(obj).find("ul").eq(0).find(">li").each(function () {
                $(this).find("a").each(function () {
                    var v = $(this).css("padding-left");
                    var i = parseInt(v) + 20;
                    $(this).css("padding-left", i + "px");
                });
                if ($(this).find(">ul").length > 0) {
                    indentTreeByLi(this)
                }
            })
        }
    }
    $("#side-markdown-view").find("a").prepend('<i class="empty-i"></i>');
    $("#side-markdown-view").find("ul").eq(0).find(">li").each(function(){
        indentTreeByLi(this)
    });
    // 绑定li点击事件
    $("#side-markdown-view").find("li").click(function(){
        if ($(this).find("ul").length > 0) {
            $(this).find("ul").eq(0).toggle();
            // 改变图标
            if ($(this).find("ul").eq(0).css('display') == "none") {
                $(this).find("i").eq(0).attr("class", "am-icon-caret-right");
            } else {
                $(this).find("i").eq(0).attr("class", "am-icon-caret-down");
            }
        } else {
            loadMarkdown($(this).find("a").eq(0).attr('href'), true)
        }
        return false
    });
    // 先将所有带子级的li关闭
    $("#side-markdown-view").find("ul").find(">li").each(function () {
        if ($(this).find("ul").length > 0) {
            $(this).trigger("click")
        }
    });
    // 默认将第一级菜单展开
    // $("#side-markdown-view").find("ul").eq(0).find(">li").each(function () {
    //     if ($(this).find("ul").length > 0) {
    //         $(this).trigger("click")
    //     }
    // });
    // 高亮并展开当前打开的地址
    highlightLiByUri(window.location.pathname);

    $("#process-mask").hide();

    // 菜单关闭隐藏
    $("#menu-icon").click(function () {
        if ($("#side-markdown-view").css("display") == "none") {
            $(this).css("left", "340px");
            $("#side-markdown-view").show();
        } else {
            $(this).css("left", "20px");
            $("#side-markdown-view").hide();
        }
    });

    // 搜索按钮
    $("#search-input button").click(function () {
        var key = $("#search-key").val();
        if (key.length == 0) {
            $("#side-menus").find("li").show();
            $("#clear-button button").hide();
        } else {
            $("#side-menus").find("li").hide();
            $("#search-loading").show();
            $.ajax({
                type     : "get",
                url      : "/search",
                data     : "key=" + encodeURIComponent(key),
                dataType : "json",
                success: function(result){
                    if (result.code == 1) {
                        for (var i = 0; i < result.data.length; i++) {
                            $("a[href='"+ result.data[i] +"']").parents("li").show();
                        }
                    }
                    $("#search-loading").hide();
                    $("#clear-button button").show();
                }
            });
        }
    });
    // 回车按钮触发搜索按钮点击事件
    $("#search-key").on("keydown", function (event) {
        if (event.keyCode == 13) {
            $("#search-input button").trigger("click");
        }
    });
    // 搜索清除按钮
    $("#clear-button button").click(function () {
        $("#search-key").val("");
        $("#search-input button").trigger("click");
    });

    
});
