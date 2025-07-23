package main

var tmpl = `
{{if .}}
<div class="bg-gray-50 rounded-xl p-4 space-y-3 max-h-60 overflow-y-auto">
    <h4 class="font-semibold text-gray-700 mb-3 flex items-center">
        <i class="fas fa-code-branch mr-2 text-indigo-500"></i>
        Select Repositories ({{len .}} found)
    </h4>
{{range .}}
    <div class="w-full">
        <label class="flex items-start p-3 w-full rounded-lg border-2 border-gray-200 hover:border-indigo-300 hover:bg-indigo-50 transition-all duration-200 cursor-pointer min-h-[80px]">
            <input type="checkbox" name="repos" value="{{.URL}}" class="mr-3 mt-1 w-4 h-4 accent-indigo-500 flex-shrink-0">
            <div class="flex-1 min-w-0">
                <div class="font-medium text-gray-800 truncate">{{.Name}}</div>
                {{if .Description}}
                    <div class="text-sm text-gray-500 mt-1 line-clamp-2">{{.Description}}</div>
                {{else}}
                    <div class="text-sm text-gray-400 mt-1 italic">No description available</div>
                {{end}}
                <div class="flex items-center mt-2 space-x-4 text-xs text-gray-400">
                    {{if .Language}}
                        <span class="flex items-center">
                            <i class="fas fa-code mr-1"></i>{{.Language}}
                        </span>
                    {{end}}
                    <span class="flex items-center">
                        <i class="fas fa-star mr-1"></i>{{.StargazersCount}}
                    </span>
                </div>
            </div>
        </label>
    </div>
{{end}}

</div>
{{else}}
<div class="text-center py-8 text-gray-500 bg-gray-50 rounded-xl">
    <i class="fas fa-search text-4xl mb-4"></i>
    <p class="text-lg font-medium">No repositories found</p>
    <p class="text-sm">This user may not exist, or has no public non-forked repositories.</p>
</div>
{{end}}
`

