package auth

import (
	"strconv"

	"github.com/xwinie/glue/core"
)

func findRoleCountByPageService() int64 {
	num, err := roleCountByPage()
	if err != nil {
		return 0
	}
	return num
}

func findRoleByPageService(p *core.Paginator) (responseEntity core.ResponseEntity) {
	roles, err := roleByPage(p.PerPageNums, p.Offset())
	var hateoas core.HateoasTemplate
	var links core.Links
	links.Add(core.LinkTo("/v1/role/{code}", "self", "GET", "根据编码获取角色信息"))
	links.Add(core.LinkTo("/v1/role/{id}", "self", "DELETE", "根据id删除角色信息"))
	links.Add(core.LinkTo("/v1/role/{id}", "self", "PUT", "根据id修改角色信息"))
	links.Add(core.LinkTo(p.PageLinkFirst(), "first", "GET", ""))
	links.Add(core.LinkTo(p.PageLinkLast(), "last", "GET", ""))
	if p.HasNext() {
		links.Add(core.LinkTo(p.PageLinkNext(), "next", "GET", ""))
	}
	if p.HasPrev() {
		links.Add(core.LinkTo(p.PageLinkPrev(), "prev", "GET", ""))
	}
	hateoas.AddLinks(links)
	type data struct {
		Roles []*SysRole
		Total int64
		*core.HateoasTemplate
	}
	d := &data{roles, p.Nums(), &hateoas}
	if err != nil {
		return *responseEntity.BuildError(core.BuildEntity(QueryError, getMsg(QueryError)))
	}
	return *responseEntity.Build(d)
}

func createRole(u SysRole) (responseEntity core.ResponseEntity) {
	if u.Code == "" {
		return *responseEntity.BuildError(core.BuildEntity(ParameterError, getMsg(ParameterError)))
	}
	isExist := u.codeIsExist()
	if isExist.Code != 100000 {
		return *responseEntity.BuildError(core.BuildEntity(RoleIsExist, getMsg(RoleIsExist)))
	}
	G, _ := core.NewGUID(1)
	id, _ := G.NextID()
	u.ID = strconv.FormatInt(id, 10)
	err := u.insert()
	if err != nil {
		return *responseEntity.BuildError(core.BuildEntity(CreateRoleError, getMsg(CreateRoleError)))
	}
	var hateoas core.Hateoas
	var links core.Links
	links.Add(core.LinkTo("/v1/role/"+u.Code, "self", "GET", "根据编码获取角色信息"))
	links.Add(core.LinkTo("/v1/role/"+strconv.FormatInt(id, 10), "self", "DELETE", "根据id删除角色信息"))
	links.Add(core.LinkTo("/v1/role/"+strconv.FormatInt(id, 10), "self", "PUT", "根据id修改角色信息"))
	hateoas.AddLinks(links)
	type data struct {
		*core.Hateoas
	}
	d := &data{&hateoas}
	return *responseEntity.BuildPostAndPut(d)
}

func deleteRoleService(id string) (responseEntity core.ResponseEntity) {
	err := deleteRole(id)
	if err != nil {
		return *responseEntity.BuildError(core.BuildEntity(DeleteRoleError, getMsg(DeleteRoleError)))
	}
	return *responseEntity.BuildDelete(core.BuildEntity(Success, getMsg(Success)))
}

func updateRoleService(id string, m SysRole) (responseEntity core.ResponseEntity) {

	m.ID = id
	m.Code = ""
	err := m.update()
	if err != nil {
		return *responseEntity.BuildError(core.BuildEntity(UpdateRoleError, getMsg(UpdateRoleError)))
	}
	u, _ := findRoleByID(id)
	var hateoas core.Hateoas
	var links core.Links
	links.Add(core.LinkTo("/v1/role/"+u.Code, "self", "GET", "根据编码获取角色信息"))
	links.Add(core.LinkTo("/v1/role/"+id, "self", "DELETE", "根据id删除角色信息"))
	links.Add(core.LinkTo("/v1/role/"+id, "self", "PUT", "根据id修改角色信息"))
	hateoas.AddLinks(links)
	type data struct {
		*core.Hateoas
	}
	d := &data{&hateoas}
	return *responseEntity.BuildPostAndPut(d)
}

func findRoleByCodeService(code string) (responseEntity core.ResponseEntity) {
	u, err := findRoleByCode(code)
	if err != nil {
		return *responseEntity.BuildError(core.BuildEntity(QueryError, getMsg(QueryError)))
	}
	var hateoas core.Hateoas
	var links core.Links
	links.Add(core.LinkTo("/v1/role/"+u.ID, "self", "DELETE", "根据id删除角色信息"))
	links.Add(core.LinkTo("/v1/role/"+u.ID, "self", "PUT", "根据id修改角色信息"))
	hateoas.AddLinks(links)
	type data struct {
		Roles SysRole
		*core.Hateoas
	}
	d := &data{u, &hateoas}
	return *responseEntity.Build(d)
}
