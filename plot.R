# read data
filepath <- "./results.csv"
all <- read.csv(filepath, header = TRUE)
all$dTotal <- all$dTotal / 1000000000

# filter by type
ints <- all[all$type=="int", ]
floats <- all[all$type=="float", ]
strings <- all[all$type=="string", ]
objects <- all[all$type=="object", ]

# color sets
degrade <- c("deepskyblue4","deepskyblue3", "deepskyblue2", "deepskyblue1")
degmix <- c("deepskyblue2","deepskyblue4", "deepskyblue1", "deepskyblue3")
accessible <- c("deepskyblue4","darkgreen", "darkorchid4", "firebrick")
colorset <- degmix

# pdf settings
wid = 8
hei = 6
fontsize = 12

# labels
lTime = "Total Time (in nanoseconds)"
lTSize = "Total Size (in bytes)"
lFormat = "Data Format"

# plot
resultsDir = "./plots/"

## total time
pdf(file=paste(resultsDir,"all-time-plot.pdf", sep="/"), height=hei, width=wid, pointsize = fontsize)
boxplot(dTotal~format, data=all, col=(colorset),
        main="Total Time",
        xlab=lFormat,
        ylab=lTime
)
dev.off()
pdf(file=paste(resultsDir,"int-time-plot.pdf", sep="/"), height=hei, width=wid, pointsize = fontsize)
boxplot(dTotal~format, data=ints, col=(colorset),
        main="64-bit Integer Total Time",
        xlab=lFormat,
        ylab=lTime
)
dev.off()
pdf(file=paste(resultsDir,"float-time-plot.pdf", sep="/"), height=hei, width=wid, pointsize = fontsize)
boxplot(dTotal~format, data=floats, col=(colorset),
        main="64-bit Floating Point Total Time",
        xlab=lFormat,
        ylab=lTime
)
dev.off()
pdf(file=paste(resultsDir,"string-time-plot.pdf", sep="/"), height=hei, width=wid, pointsize = fontsize)
boxplot(dTotal~format, data=strings, col=(colorset),
        main="Strings Total Time",
        xlab=lFormat,
        ylab=lTime
)
dev.off()
pdf(file=paste(resultsDir,"object-time-plot.pdf", sep="/"), height=hei, width=wid, pointsize = fontsize)
boxplot(dTotal~format, data=objects, col=(colorset),
        main="Complex Object Structures Total Time",
        xlab=lFormat,
        ylab=lTime
)
dev.off()