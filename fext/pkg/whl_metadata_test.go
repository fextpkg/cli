package pkg

var (
	PackageDependencies  = []string{"charset-normalizer", "idna", "urllib3", "certifi"}
	PackageExtras        = []string{"socks", "use_chardet_on_py3"}
	PackageExtrasUnknown = []string{"test", "another"}
)

const (
	PackageName    = "test_requests"
	PackageVersion = "2.31.0"

	Metadata = `Metadata-Version: 2.1
Name: test-requests
Version: 2.31.0
Summary: Python HTTP for Humans.
Home-page: https://requests.readthedocs.io
Author: Kenneth Reitz
Author-email: me@kennethreitz.org
License: Apache 2.0
Project-URL: Documentation, https://requests.readthedocs.io
Project-URL: Source, https://github.com/psf/requests
Platform: UNKNOWN
Classifier: Development Status :: 5 - Production/Stable
Classifier: Environment :: Web Environment
Classifier: Intended Audience :: Developers
Classifier: License :: OSI Approved :: Apache Software License
Classifier: Natural Language :: English
Classifier: Operating System :: OS Independent
Classifier: Programming Language :: Python
Classifier: Programming Language :: Python :: 3
Classifier: Programming Language :: Python :: 3.7
Classifier: Programming Language :: Python :: 3.8
Classifier: Programming Language :: Python :: 3.9
Classifier: Programming Language :: Python :: 3.10
Classifier: Programming Language :: Python :: 3.11
Classifier: Programming Language :: Python :: 3 :: Only
Classifier: Programming Language :: Python :: Implementation :: CPython
Classifier: Programming Language :: Python :: Implementation :: PyPy
Classifier: Topic :: Internet :: WWW/HTTP
Classifier: Topic :: Software Development :: Libraries
Requires-Python: >=3.7
Description-Content-Type: text/markdown
License-File: LICENSE
Requires-Dist: charset-normalizer (<4,>=2)
Requires-Dist: idna (<4,>=2.5)
Requires-Dist: urllib3 (<3,>=1.21.1)
Requires-Dist: certifi (>=2017.4.17)
Provides-Extra: security
Provides-Extra: socks
Requires-Dist: PySocks (!=1.5.7,>=1.5.6) ; extra == 'socks'
Provides-Extra: use_chardet_on_py3
Requires-Dist: chardet (<6,>=3.0.2) ; extra == 'use_chardet_on_py3'
`
	MetadataBrokenMarker = `Metadata-Version: 2.1
Name: test-requests
Version: 2.31.0
Requires-Dist: charset-normalizer (<4,>=2)
Requires-Dist: idna (<4,>=2.5)
Requires-Dist: urllib3 (<3,>=1.21.1)
Requires-Dist: certifi (>=2017.4.17)
Provides-Extra: security
Provides-Extra: socks
Requires-Dist: PySocks (!=1.5.7,>=1.5.6) ; unknown == 'test' and extra == 'socks'
Provides-Extra: use_chardet_on_py3
Requires-Dist: chardet (<6,>=3.0.2) ; extra == 'use_chardet_on_py3'
`
	MetadataInvalidSyntaxMarker = `Metadata-Version: 2.1
Name: test-requests
Version: 2.31.0
Requires-Dist: charset-normalizer (<4,>=2)
Requires-Dist: idna (<4,>=2.5)
Requires-Dist: urllib3 (<3,>=1.21.1)
Requires-Dist: certifi (>=2017.4.17)
Provides-Extra: security
Provides-Extra: socks
Requires-Dist: PySocks (!=1.5.7,>=1.5.6) ; == 'test' extra == 'socks'
Provides-Extra: use_chardet_on_py3
Requires-Dist: chardet (<6,>=3.0.2) ; extra == 'use_chardet_on_py3'
`
)
