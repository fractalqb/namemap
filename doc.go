// Package namemap addresses the situation where a software has to
// mediate between different data sources where each source (domain)
// has a different ID (name) for one thing. The NameMap provides an
// efficient tool to map such names from one domain to another.
//
// Besides the types and methods for name mapping this package
// provides a simple file format to define such mappings that can
// easily be parsed for application.
// As it is often desirable one of the domains can be defined as
// the standard domain the file format allows for easily define one
// domain as the standard domain. E.g.
//
//   [\input output l10n:EN l10n:DE]
//   ( note  rem    remark  Bemerkung)
//   ( warn  warnig warning Warnung)
//
// defines a map that has an 'input', 'output' domain an additionally
// has some natural language domains for user interfaces
// (e.g. GUI). The 'input' domain is selected to be the standard domain.
package namemap

//go:generate versioner -pkg namemap -bno build_no ./VERSION ./version.go
