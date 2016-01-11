package main

import (
    "fmt"
    "bytes"
    "os"
    "flag"
    "strings"
    "strconv"
    "encoding/xml"
    "io/ioutil"

)


type dependency struct{
  Name string `xml:"id,attr"`
  Version string `xml:"version,attr"`
}

type file struct {
  Source string `xml:"src,attr"`
  Target string `xml:"target,attr"`
}

type nuspec struct {
  XMLName     string `xml:"package,InXMLNameTag"`
  Namespace   string `xml:"xmlns,attr"`
  Id          string `xml:"metadata>id"`
  Title       string `xml:"metadata>title"`
  Version     string `xml:"metadata>version"`
  Authors     string `xml:"metadata>authors"`
  Owners      string `xml:"metadata>owners"`
  Description string `xml:"metadata>description"`
  Summary     string `xml:"metadata>summary"`
  ProjectUrl  string `xml:"metadata>projectUrl"`
  Tags        string `xml:"metadata>tags"`
  Copyright   string `xml:"metadata>copyright"`
  LicenseUrl  string `xml:"metadata>licenseUrl"`
  rLA         bool   `xml:"metadata>requireLicenseAcceptance"`
  Deps        []dependency `xml:"metadata>dependencies>dependency"`
  RelNotes    string `xml:"metadata>releaseNotes"` 
  Files       []file `xml:"files>file"`
  
} 

type multiflag []string

func (d *multiflag) String() string {
    return fmt.Sprintf("%d", *d)
}
 
func (d *multiflag) Set(value string) error {
    *d = append(*d, value)
    return nil
}

var alldeps multiflag
var allfiles multiflag
var allcommands multiflag

func main() {


  cmds := map[string] string{"ipkg":"Install-ChocolateyPackage", "izip":"Install-ChocolateyZipPackage", "uzip":"Get-ChocolateyUnzip" }


  flag.Var(&alldeps, "deps", "Depedencies and versions -deps dep1,ver1 -deps dep2,ver2")
  flag.Var(&allfiles, "files", "Files source and target -files file_src1,file_trg1 -files file_src2,file_trg2")
  flag.Var(&allcommands, "commands", "All PS commands you want to execute in the form command:arg1,arg2...,argN -commands ipkg,url,target,quiet commands izip,url,target,quiet")
  versionPtr := flag.String("version", "", "version of the package to build in the format x.y.z")
  namePtr    := flag.String("name", "", "Id,name,title of the package to build")
  descPtr    := flag.String("desc", "My custom package", "Description of the package to build")
  sumPtr     := flag.String("sum", "My custom package", "Summary of the package to build")
  alcPtr     := flag.Bool("license", false, "Does the user need to accept the license?")
  authorsPtr := flag.String("authors", "Me", "Authors of the packaged software")
  ownersPtr  := flag.String("owners", "Me", "Owners of the packaged software")
  lUrlPtr    := flag.String("lurl", "http://github.com/me/my_project/license", "License URL of the packaged software")
  pUrlPtr    := flag.String("purl", "http://github.com/me/my_project", "Project URL of the packaged software")

  flag.Parse()

  
  if ( (len(*namePtr) <= 0 || len(*versionPtr) <= 0) ) {
    fmt.Println("Please provide the following required parameters:")
    flag.PrintDefaults()
    return
  }


  v := &nuspec{Id: *namePtr, Title: *namePtr, Namespace: "http://schemas.microsoft.com/packaging/2010/07/nuspec.xsd", Version: *versionPtr}
  v.Authors = *authorsPtr
  v.Owners = *ownersPtr
  v.Description = *descPtr
  v.Summary = *sumPtr
  v.ProjectUrl = *pUrlPtr
  v.LicenseUrl = *lUrlPtr
  v.rLA = *alcPtr


  for _,d := range alldeps {
    tmp := strings.Split(d, ",")
    v.Deps = append(v.Deps, dependency{Name: tmp[0], Version: tmp[1]})
  }
  for _,f := range allfiles {
    tmp := strings.Split(f, ",")
    v.Files = append(v.Files, file{Source: tmp[0], Target: tmp[1]})
  }
  

  cnt_commands := len(allcommands)

  if(cnt_commands>0){
    
    ordered :=  make(map[int] bytes.Buffer, cnt_commands)
    for _,c := range allcommands {
      var bstring bytes.Buffer
      
      tmp := strings.Split(c, ",")
      if(tmp[1] != "custom") {
        bstring.WriteString(cmds[tmp[1]])
        bstring.WriteString(" ")
      }     
      
      tmp2 := tmp[2:]
      for _,v := range tmp2 {
        bstring.WriteString(v)
        bstring.WriteString(" ")
      }
      
      bstring.WriteString("\r\n")
      idx, c_err := strconv.Atoi(tmp[0])
      if c_err != nil {
        fmt.Print("Command index needs to be an integer from 1 to N: ")
        panic(c_err)
      }
      ordered[idx] = bstring
      
    }

    
    

    psfile, ioerr := os.OpenFile("ChocolateyInstall.ps1", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)

    for i:=1; i<=len(allcommands); i++ {
      tmp3 := ordered[i]
      _, ioerr = psfile.Write(tmp3.Bytes())
    }
  
    if ioerr != nil {
      panic(ioerr)
    }

    v.Files = append(v.Files, file{Source: "ChocolateyInstall.ps1", Target: "tools"})
    
  }

  output, err := xml.MarshalIndent(v, "  ", "    ")
  if err != nil {
    fmt.Printf("error: %v\n", err)
  }
  

  xml_header := []byte(xml.Header)
  xml_out := make([]byte, len(output)+len(xml_header))
  cbytes := copy(xml_out[0:], xml_header)
  copy(xml_out[cbytes:], output)

  ioerr2 := ioutil.WriteFile("package.nuspec", xml_out, 0644)
  
  if ioerr2 != nil {
    panic(ioerr2)
  }

}


