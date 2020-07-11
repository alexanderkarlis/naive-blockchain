import React, { useState } from 'react';
import { Theme, createStyles, makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import { TextField, Typography } from '@material-ui/core';

const useStyles = makeStyles((theme: Theme) =>
  createStyles({
    root: {
      width: "100%",
      display: 'flex',
      flexWrap: 'wrap',
      '& > *': {
        margin: theme.spacing(1),
        width: theme.spacing(16),
        height: theme.spacing(16),
      },
    },
    paper: {
      width: "100%"
    },
    textFormat: {
      fontSize: "14px"
    }
  }),
);


export const BlockComponent = (props) => {
    const classes = useStyles();

    return (
        <div className={classes.root}>
            <Paper className={classes.paper} elevation={3}>
              <div>
                <Typography className={classes.textFormat} variant='h4'>{`ID: ${props.block.ID}`}</Typography>
              </div>
              <div>
                <Typography className={classes.textFormat} variant='h6'>{`Index: ${props.block.Index}`}</Typography>
              </div>
              <div>
                <Typography className={classes.textFormat} variant='h4'>{`Hash: ${props.block.Hash}`}</Typography>
              </div>
              <div>
                <Typography className={classes.textFormat} variant='h4'>{`PreviousHash: ${props.block.PreviousHash}`}</Typography>
              </div>
              <div>
                <Typography className={classes.textFormat} variant='h4'>{`Timestamp: ${props.block.Timestamp}`}</Typography>
              </div>
              <div>
                <Typography className={classes.textFormat} variant='h4'>{`Data: ${props.block.Data}`}</Typography>
              </div>
            </Paper>
        </div>
    )
}