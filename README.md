## Overview

A simple Go helper utility that creates a nuspec file and a chocolatey powershell script to assist with the automated build of a custom local chocolatey package.      


## Setup


### Setup Requirements

No special requirements. 
You only need to run the tool in the workspace directory where you chocolatey package will be built.
Make sure you appropriately escape all special chars in the parameters for your shell/cli.

### Beginning with makechoco

## Usage

The tool accepts the following parameters:

```
 -authors string
      Authors of the packaged software (default "Me")
  -commands value
      All PS commands you want to execute in the form order_no,command:arg1,arg2...,argN -commands 1,ipkg,url,target,quiet -commands 2,izip,url,target,quiet (default []) ( command types: ipkg,izip,uzip,custom order of execution needs to be sequential and starts at 1 e.g. 1,2,3...N )
  -deps value
      Depedencies and versions -deps dep1,ver1 -deps dep2,ver2 (default [])
  -desc string
      Description of the package to build (default "My custom package")
  -files value
      Files source and target -files file_src1,file_trg1 -files file_src2,file_trg2 (default [])
  -license
      Does the user need to accept the license?
  -lurl string
      License URL of the packaged software (default "http://github.com/me/my_project/license")
  -name string
      Id,name,title of the package to build
  -owners string
      Owners of the packaged software (default "Me")
  -purl string
      Project URL of the packaged software (default "http://github.com/me/my_project")
  -sum string
      Summary of the package to build (default "My custom package")
  -version string
      version of the package to build in the format x.y.z
```

Example:

```
makechoco -deps 7zip,15.14 -name mypackage -version 1.0.0 -files custom.zip,content -commands 2,uzip,\$zip_src,C:\\tmp --commands 1,custom,"\$zip_src = Join-Path \"\$env:ChocolateyPackageFolder\" 'content\custom.zip'"

```


### Limitations

Currently only the following chocolatey powershell commands are supported:
- Install-ChocolateyPackage - ipkg
- Install-ChocolateyZipPackage - izip
- Get-ChocolateyUnzip - uzip

But a custom options is allowed to use powershell statements and commands.






