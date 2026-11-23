% jrresume: the design of Jon Reyna's résumé. vitae writes the content
% (build/resume/resume.<variant>.tex); this class decides how it looks, and
% `make resume-pdf` builds each document with latexmk and pdfLaTeX.
%
% The design carries over the résumé sent in November 2026 (tag sent-2026-11 in the
% old résumé repository): a single column in reading order, so text extractors and
% applicant-tracking parsers read it top to bottom; real text throughout (no tables,
% text boxes, icons or images); Lato with f-ligatures off; no hyphenation; one deep
% blue for structure, near-black text and gray for metadata.
%
% What changed for generated content: no hand-placed page break. Instead no paragraph
% (so no bullet) splits across pages, and a job heading never ends a page.
\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{jrresume}[2026/11/22 Jon Reyna's résumé]
\LoadClass[10pt,letterpaper]{article}

\RequirePackage[T1]{fontenc}
\RequirePackage[utf8]{inputenc}
% Lato extracts cleanly; f-ligatures are off so parsers see "ff", "fi" and "ft"
% as plain letters (Source Sans Pro's ligatures extracted as "Sotfware").
\RequirePackage[default,lining]{lato}
\RequirePackage{microtype}
\DisableLigatures[f]{encoding = *, family = *}
\RequirePackage[letterpaper,left=0.7in,right=0.7in,top=0.55in,bottom=0.55in]{geometry}
\RequirePackage{enumitem}
\RequirePackage{titlesec}
\RequirePackage{xcolor}
\RequirePackage{needspace}
\RequirePackage[hidelinks,pdfusetitle]{hyperref}

% Map glyphs to Unicode so copied and extracted text matches what is printed.
\input{glyphtounicode}
\pdfgentounicode=1
% Emit real space characters, so extractors don't run words together.
\pdfinterwordspaceon

% Palette: one deep blue for structure, near-black text, gray for metadata.
\definecolor{ink}{HTML}{222222}
\definecolor{accent}{HTML}{1F4E79}
\definecolor{muted}{HTML}{666666}
\definecolor{rule}{HTML}{A7B6C6}

\pagestyle{empty}
\AtBeginDocument{\color{ink}}
\linespread{1.06}
% No hyphenation: a word split across lines reaches a parser as two fragments.
\hyphenpenalty=10000
\exhyphenpenalty=10000
\emergencystretch=2em
% Ragged right: without hyphenation, justified lines open wide gaps.
\AtBeginDocument{\raggedright}
\setlength{\parindent}{0pt}
\setlength{\parskip}{0pt}
% Never break a page inside a paragraph, so a bullet always stays whole.
\interlinepenalty=10000

% Weights: Lato Black for the name, Semibold for structure, Regular for text.
\newcommand{\semibold}{\fontseries{sb}\selectfont}
\newcommand{\black}{\fontseries{eb}\selectfont}

% \sep{}: the middle dot between items of a line; the line may break after it.
\newcommand{\sep}{\enspace\textperiodcentered\enspace\allowbreak}

% \resumeheader{Name}{Label}{Contact line}
\newcommand{\resumeheader}[3]{%
  {\fontsize{26}{28}\black\color{accent}#1}\par
  \vspace{5pt}%
  {\large\semibold #2}\par
  \vspace{4pt}%
  {\small\color{muted}#3}\par
  \vspace{8pt}%
  {\color{accent}\rule{\linewidth}{1.2pt}}\par
  \vspace{4pt}}

% Section heads: semibold, accent color, a hairline rule running to the margin.
\titleformat{\section}
  {\bfseries\color{accent}\large}{}{0pt}{}
  [\vspace{-4pt}{\color{rule}\titlerule[0.6pt]}]
\titlespacing*{\section}{0pt}{11pt}{5pt}

% Highlights: a list of bullets under a job.
\newlist{highlights}{itemize}{1}
\setlist[highlights]{leftmargin=1.05em,labelsep=0.5em,itemsep=2.5pt,parsep=0pt,
  topsep=3pt,label={\color{muted}\small\textbullet}}

% \job{Company}{Title}{Location}{Dates}: company and dates on the first line,
% title and location on the second, so extraction keeps the order. The first job
% after a section head sits closer to it.
\newif\ifjr@afterhead
\newcommand{\job}[4]{%
  \par\ifjr@afterhead\vspace{2pt}\else\vspace{8pt}\fi
  \global\jr@afterheadfalse
  \needspace{4\baselineskip}%
  {\bfseries\large #1}\hfill{\color{muted}#4}\par\nopagebreak
  \vspace{0.5pt}\nopagebreak
  {\color{accent}#2}{\color{muted}\enspace\textperiodcentered\enspace #3}\par\nopagebreak}
\let\jr@section\section
% A section head never ends a page: it needs room for itself and a few lines.
\renewcommand{\section}{\needspace{6\baselineskip}\global\jr@afterheadtrue\jr@section}

% \lead{Phrase}: a semibold lead-in for a bullet, for the skim.
\newcommand{\lead}[1]{{\semibold #1.}\ }

% \skill{Label}{List}: labels in a fixed column; wrapped lines hang under the list.
\newlength{\skilllabel}
\setlength{\skilllabel}{11.4em}
\newcommand{\skill}[2]{%
  \par\hangindent=\skilllabel\hangafter=1
  \makebox[\skilllabel][l]{\semibold #1}#2\par\vspace{2pt}}

% \earlier{Text}: the one-line note on earlier work, after the last job.
\newcommand{\earlier}[1]{%
  \par\vspace{6pt}{\color{muted}\textit{Earlier:} #1}\par}

% \school{School}{Years}{Program}{Minors}: like a job's two lines; Minors may be empty.
% The paragraph that follows is the note, standing and coursework.
\newcommand{\school}[4]{%
  \par\ifjr@afterhead\else\vspace{8pt}\fi
  \global\jr@afterheadfalse
  \needspace{4\baselineskip}%
  {\bfseries\large #1}\hfill{\color{muted}#2}\par\nopagebreak
  \vspace{0.5pt}\nopagebreak
  {\color{accent}#3}%
  \if\relax\detokenize{#4}\relax\else{\color{muted}\enspace\textperiodcentered\enspace #4}\fi
  \par\nopagebreak
  \vspace{3pt}\nopagebreak}

% \maxpages{N}: the variant's page limit (max_pages in resume.yaml). A document that
% runs past it stops the build, so a too-long version is never published.
\newcommand{\jr@maxpages}{0}
\newcommand{\maxpages}[1]{\renewcommand{\jr@maxpages}{#1}}
\AtEndDocument{%
  \ifnum\jr@maxpages>0
    \ifnum\value{page}>\jr@maxpages\relax
      \ClassError{jrresume}{This version runs to \number\value{page} pages,
        but max_pages allows \jr@maxpages}{Trim the variant in resume/resume.yaml.}%
    \fi
  \fi}
