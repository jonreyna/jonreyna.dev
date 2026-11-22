package variant

import (
	"fmt"
	"slices"
	"sort"

	"github.com/jonreyna/jonreyna.dev/tools/vitae/internal/resume"
)

// DefaultBoost is what each boosted tag adds to a highlight's score when the variant
// doesn't set BoostBy.
const DefaultBoost = 20

// Apply returns r with the named variant applied. It doesn't validate r; run package
// lint first. The error wraps [ErrUnknown] if r has no such variant.
func Apply(r *resume.Resume, name string) (*Resolved, error) {
	v, ok := r.Variants[name]
	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrUnknown, name)
	}
	out := &Resolved{Variant: name, Sections: v.Sections, Basics: r.Basics}
	if v.Label != "" {
		out.Basics.Label = v.Label
	}
	if out.Has(resume.SectionSummary) {
		key := v.Summary
		if key == "" {
			key = "default"
		}
		out.Summary = r.Summary[key]
	}
	if out.Has(resume.SectionExperience) {
		out.Experience = roles(r.Experience, &v)
	}
	if out.Has(resume.SectionEarlier) {
		out.Earlier = r.Earlier
	}
	if out.Has(resume.SectionSkills) {
		out.Skills = r.Skills
	}
	if out.Has(resume.SectionEducation) {
		out.Education = make([]resume.School, 0, len(r.Education))
		for i := range r.Education {
			s := r.Education[i]
			s.Courses = courses(s.Courses, v.ShowCourses)
			out.Education = append(out.Education, s)
		}
	}
	return out, nil
}

// roles keeps the roles v selects, each with its selected highlights.
func roles(all []resume.Role, v *resume.Variant) []resume.Role {
	out := make([]resume.Role, 0, len(all))
	for _, role := range all {
		if v.Since != nil && role.End != nil && role.End.Before(*v.Since) {
			continue
		}
		if v.MaxRoles > 0 && len(out) == v.MaxRoles {
			break
		}
		role.Highlights = highlights(role.Highlights, v)
		out = append(out, role)
	}
	return out
}

// highlights keeps the highlights v selects, in their original order.
func highlights(all []resume.Highlight, v *resume.Variant) []resume.Highlight {
	boostBy := v.BoostBy
	if boostBy == 0 {
		boostBy = DefaultBoost
	}
	kept := make([]scored, 0, len(all))
	for i, h := range all {
		if anyTag(h.Tags, v.Exclude) {
			continue
		}
		score := h.Weight
		for _, t := range h.Tags {
			if slices.Contains(v.Boost, t) {
				score += boostBy
			}
		}
		if score >= v.MinWeight {
			kept = append(kept, scored{index: i, score: score})
		}
	}
	if v.MaxHighlights > 0 && len(kept) > v.MaxHighlights {
		// Highest score first; the stable sort leaves ties in file order.
		sort.SliceStable(kept, func(a, b int) bool { return kept[a].score > kept[b].score })
		kept = kept[:v.MaxHighlights]
		sort.Slice(kept, func(a, b int) bool { return kept[a].index < kept[b].index })
	}
	out := make([]resume.Highlight, 0, len(kept))
	for _, k := range kept {
		out = append(out, all[k.index])
	}
	return out
}

// courses keeps the courses carrying a tag in show, or all of them if show is empty.
func courses(all []resume.Course, show []string) []resume.Course {
	if len(show) == 0 {
		return all
	}
	out := make([]resume.Course, 0, len(all))
	for _, c := range all {
		if anyTag(c.Tags, show) {
			out = append(out, c)
		}
	}
	return out
}

// anyTag reports whether tags and set share a tag.
func anyTag(tags, set []string) bool {
	return slices.ContainsFunc(tags, func(t string) bool { return slices.Contains(set, t) })
}
