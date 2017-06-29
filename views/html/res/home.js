function HomePage(datagrid_id, dialog_id) {
    this.datagrid_id = datagrid_id;
    this.dialog_id = dialog_id;

    this.dialog = $('#' + this.dialog_id).dialog('center');
}

HomePage.prototype = {
    //原型字面量方式会将对象的constructor变为Object，此外强制指回ResHome
    constructor: HomePage,
    baseUrl: '/res/',
    loadDatagrid: function () {
        this.datagrid = $('#' + this.datagrid_id).datagrid({
            url: '/res?pid=-1',
            pagination: true,
            fitColumns: true,
            rownumbers: true,
            fit: true,
            toolbar: '#toolbar',
            method: 'get',
            singleSelect: true,
            columns: [[
                {
                    field: 'id', title: 'Id', sortable: true
                }, {
                    field: 'res_name', title: '资源名称', width: 100, sortable: true
                }, {
                    field: 'res_type', title: '资源类型', width: 100, sortable: true,
                    formatter: function (value) {
                        if (value == 0) {
                            return '菜单'
                        } else {
                            return '资源'
                        }
                    }
                }, {
                    field: 'path', title: '路径', width: 100, sortable: true
                }, {
                    field: 'create_time', title: '创建时间', width: 100, sortable: true,
                    formatter: function (value, row, index) {
                        return new Date(value).format("yyyy-MM-dd hh:mm:ss");
                    }
                }, {
                    field: 'update_time', title: '最后一次更新时间', width: 100, sortable: true,
                    formatter: function (value, row, index) {
                        return new Date(value).format("yyyy-MM-dd hh:mm:ss");
                    }
                }
            ]],
            onSelect: function (rowIndex, rowData) {
                $('table.ddv.datagrid-f').each(function (index, ddv) {
                    $(ddv).datagrid('unselectAll');
                });
                console.info(this)
            },
            view: detailview,
            detailFormatter: function (index, row) {
                return '<div style="padding:2px"><table class="ddv"></table></div>';
            },
            onExpandRow: function (index, row) {
                var ddv = $(this).datagrid('getRowDetail', index).find('table.ddv');
                ddv.datagrid({
                    method: 'get',
                    url: '/res?pid=' + row.id,
                    fitColumns: true,
                    singleSelect: true,
                    loadMsg: '',
                    height: 'auto',
                    columns: [[
                        {field: 'id', title: 'Id', sortable: true},
                        {
                            field: 'res_name', title: '角色名', width: 100, sortable: true
                        }, {
                            field: 'res_type', title: '资源类型', width: 100, sortable: true,
                            formatter: function (value) {
                                if (value == 0) {
                                    return '菜单'
                                } else {
                                    return '资源'
                                }
                            }
                        }, {
                            field: 'path', title: '路径', width: 100, sortable: true
                        }, {
                            field: 'create_time', title: '创建时间', width: 100, sortable: true,
                            formatter: function (value, row, index) {
                                return new Date(value).format("yyyy-MM-dd hh:mm:ss");
                            }
                        }, {
                            field: 'update_time', title: '最后一次更新时间', width: 100, sortable: true,
                            formatter: function (value, row, index) {
                                return new Date(value).format("yyyy-MM-dd hh:mm:ss");
                            }
                        }
                    ]],
                    onResize: function () {
                        $('#dg').datagrid('fixDetailRowHeight', index);
                    },
                    onLoadSuccess: function () {
                        setTimeout(function () {
                            $('#dg').datagrid('fixDetailRowHeight', index);
                        }, 0);
                    },
                    onSelect: function (rowIndex, rowData) {
                        var cur = this;
                        $('table.ddv.datagrid-f').each(function (index, item) {
                            if (cur != item) {
                                $(item).datagrid('unselectAll');
                            }
                        });
                        $(parent).datagrid('unselectAll');
                        // console.info($(this).datagrid('getParentGrid'))
                        // var ddv = $(this).datagrid('getRowDetail', index).find('table.ddv');
                        // console.info($(this).datagrid('getParentGrid'))
                    },
                });
                $('#dg').datagrid('fixDetailRowHeight', index);
            }
        });
        return this
    },

    query: function (queryParams) {
        console.debug(queryParams);
        $('#' + this.datagrid_id).datagrid({
            queryParams: queryParams
        });
        return this;
    },

    add: function () {
        this.dialog.dialog('open').dialog('setTitle', '添加');
        initForm();
        this.url = this.baseUrl;
    },

    edit: function () {
        var row = this.datagrid.datagrid('getSelected');
        if (row == null) {
            $('table.ddv.datagrid-f').each(function (index, item) {
                row = $(item).datagrid('getSelected');
            });
        }
        if (row) {
            initForm();
            this.dialog.dialog('open').dialog('setTitle', '编辑');
            if (row.pid == -1) {
                $('#fm').form('load', {
                    res_name: row.res_name,
                    path: row.path
                })
            } else {
                $('#fm').form('load', row);
            }
            url = this.baseUrl + row.id;
        } else {
            $.alertMsg('请选择一行数据');
        }
    },

    save: function () {
        var homne = this;
        $('#fm').form('submit', {
            url: homne.url,
            onSubmit: function () {
                if ($(this).form('validate')) {
                    $.messager.progress();
                    return true;
                }
                return false;
            },
            success: function (result) {
                $.messager.progress('close');
                var result = eval('(' + result + ')');
                if (result.ok) {
                    $('#cc').combobox('reload');
                    $('#dlg').dialog('close');        // close the dialog
                    $('#dg').datagrid('reload');    // reload the user data
                } else {
                    $.messager.show({
                        title: '系统提示',
                        msg: result.msg
                    });
                }
            }
        });
    },

    destroy: function () {
        var row = $('#dg').datagrid('getSelected');
        if (row == null) {
            $('table.ddv.datagrid-f').each(function (index, item) {
                row = $(item).datagrid('getSelected');
            });
        }
        if (row) {
            $.messager.confirm('系统提醒', '您确定删除这条数据吗?', function (r) {
                if (r) {
                    $.messager.progress();
                    $.ajax({
                        url: base + row.id,
                        type: 'DELETE',
                        dataType: 'json',
                        success: function (result) {
                            $.messager.progress('close');
                            if (result.ok) {
                                $('#dg').datagrid('reload');    // reload the user data
                            } else {
                                $.messager.show({    // show error message
                                    title: '系统提示',
                                    msg: result.msg
                                });
                            }
                        }
                    });
                }
            });
        }
    }
};

// =========================================================================================================================================================
var resHomePage
$(function () {
    resHomePage = new HomePage('dg', 'dlg')
    resHomePage.loadDatagrid();
    console.info(resHomePage)
});

function initForm() {
    $('#res_type').combobox({value: 0});
    $('#seq').numberspinner({value: 0});
    $('#path').textbox({disabled: false});
}


